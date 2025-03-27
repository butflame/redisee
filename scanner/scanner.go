package scanner

import (
	"context"
	"fmt"
	"redisee/config"
	"redisee/counter"
	"regexp"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

var allKeyTypes = []string{"string", "hash", "list", "set", "zset"}

type dbToScan struct {
	dbIndex  int
	keyCount int
	client   *redis.Client
}

type Scanner struct {
	wg          sync.WaitGroup
	dbs         []dbToScan
	finishChan  chan struct{}
	scanJobChan chan scanJob
}

func New() *Scanner {
	scanner := &Scanner{
		finishChan:  make(chan struct{}),
		scanJobChan: make(chan scanJob),
	}

	return scanner
}

type scanJob struct {
	dbClient *redis.Client
	keyType  string
	keys     []string
}

func (s *Scanner) setupClients(ctx context.Context) {
	defaultClient := redis.NewClient(&redis.Options{
		Addr:     config.Config.Host + ":" + strconv.Itoa(config.Config.Port),
		Password: config.Config.Password,
		DB:       0,
	})
	defer defaultClient.Close()

	// check if memory usage is available
	_, err := defaultClient.MemoryUsage(ctx, "").Result()
	if err != nil && redis.Nil != err {
		panic(err)
	}
	infoMap, err := defaultClient.InfoMap(ctx).Result()
	if err != nil {
		panic(err)
	}
	keySpaceMap := infoMap["Keyspace"]
	numberRe := regexp.MustCompile(`\d+`)
	for dbKey, dbInfo := range keySpaceMap {
		dbIndex, err := strconv.Atoi(numberRe.FindString(dbKey))
		if err != nil {
			panic(err)
		}
		if config.Config.Db == config.ALL_DB || dbIndex == config.Config.Db {
			keyCount, err := strconv.Atoi(numberRe.FindString(dbInfo))
			if err != nil {
				panic(err)
			}
			// 创建新的 Redis 客户端
			client := redis.NewClient(&redis.Options{
				Addr:     config.Config.Host + ":" + strconv.Itoa(config.Config.Port),
				Password: config.Config.Password,
				DB:       dbIndex,
			})
			s.dbs = append(s.dbs, dbToScan{
				dbIndex:  dbIndex,
				keyCount: keyCount,
				client:   client,
			})
		}
	}
	if len(s.dbs) == 0 {
		panic("the db you specified has no keys")
	}
}

func (s *Scanner) setOverallStats(ctx context.Context) {
	if len(s.dbs) == 0 {
		panic("no db clients, make sure you have setup clients")
	}
	client := s.dbs[0].client
	infoMap, err := client.InfoMap(ctx).Result()
	if err != nil {
		panic(err)
	}
	// 从 infoMap 中提取所需的信息
	usedMemory, _ := strconv.ParseInt(infoMap["Memory"]["used_memory"], 10, 64)
	maxMemory, _ := strconv.ParseInt(infoMap["Memory"]["maxmemory"], 10, 64)
	counter.SetTotalMemory(maxMemory)
	counter.SetUsedMemory(usedMemory)
}

func (s *Scanner) startScanWorker(ctx context.Context) {
	for job := range s.scanJobChan {
		pipeline := job.dbClient.Pipeline()
		for _, key := range job.keys {
			pipeline.TTL(ctx, key)
			pipeline.MemoryUsage(ctx, key)
		}
		cmds, err := pipeline.Exec(ctx)
		if err != nil {
			panic(err)
		}
		for i, key := range job.keys {
			ttlCmd := cmds[i*2]
			memoryUsageCmd := cmds[i*2+1]
			if ttlCmd.Err() != nil {
				if redis.Nil == ttlCmd.Err() {
					continue
				}
				fmt.Printf("error getting ttl for key %s: %v\n", key, ttlCmd.Err())
			}
			if memoryUsageCmd.Err() != nil {
				if redis.Nil == memoryUsageCmd.Err() {
					continue
				}
				fmt.Printf("error getting memory usage for key %s: %v\n", key, memoryUsageCmd.Err())
			}
			ttl := ttlCmd.(*redis.DurationCmd).Val()
			memoryUsage := memoryUsageCmd.(*redis.IntCmd).Val()
			counter.AddKey(key, job.keyType, ttl, memoryUsage)
		}
	}
	s.wg.Done()
}

func (s *Scanner) Run(ctx context.Context) {
	s.setupClients(ctx)
	fmt.Printf("Scanning %d db(s), %d concurrency, scan pattern: %s\n", len(s.dbs), config.Config.Concurrency, config.Config.ScanPattern)
	s.setOverallStats(ctx)
	for i := 0; i < config.Config.Concurrency; i++ {
		s.wg.Add(1)
		go s.startScanWorker(ctx)
	}
	for _, db := range s.dbs {
		client := db.client
		totalScanned := 0
		for _, keyType := range allKeyTypes {
			var cursor uint64 = 0
			for {
				keys, nextCursor, err := client.ScanType(ctx, cursor, config.Config.ScanPattern, 1000, keyType).Result()
				if err != nil {
					panic(err)
				}
				// Split keys into groups of at most 100 keys each
				for i := 0; i < len(keys); i += 100 {
					end := min(i+100, len(keys))
					group := keys[i:end]

					s.scanJobChan <- scanJob{
						dbClient: client,
						keyType:  keyType,
						keys:     group,
					}
				}
				cursor = nextCursor
				totalScanned += len(keys)
				// 使用 \r 回到行首并覆盖上一行内容
				if totalScanned > 0 {
					fmt.Printf("\rScanning db %d, keys %d/%d scanned", db.dbIndex, totalScanned, db.keyCount)
				}
				if cursor == 0 {
					break
				}
			}
		}
		fmt.Println()
	}
	close(s.scanJobChan)
	s.wg.Wait()

	defer fmt.Printf("Finished scanning all dbs\n\n\n")
	defer func() {
		for _, db := range s.dbs {
			db.client.Close()
		}
	}()
	defer close(s.finishChan)
}

func (s *Scanner) Finish() chan struct{} {
	return s.finishChan
}
