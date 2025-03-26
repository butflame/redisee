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

type Scanner struct {
	wg          sync.WaitGroup
	dbClients   []*redis.Client
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
	var defaultClient *redis.Client
	if config.Config.Db == config.ALL_DB {
		defaultClient = redis.NewClient(&redis.Options{
			Addr:     config.Config.Host + ":" + strconv.Itoa(config.Config.Port),
			Password: config.Config.Password,
			DB:       0,
		})
		infoMap, err := defaultClient.InfoMap(ctx).Result()
		if err != nil {
			panic(err)
		}
		keySpaceMap := infoMap["Keyspace"]
		numberRe := regexp.MustCompile(`\d+`)
		for dbKey := range keySpaceMap {
			dbIndex, err := strconv.Atoi(numberRe.FindString(dbKey))
			if err != nil {
				panic(err)
			}
			// 创建新的 Redis 客户端
			client := redis.NewClient(&redis.Options{
				Addr:     config.Config.Host + ":" + strconv.Itoa(config.Config.Port),
				Password: config.Config.Password,
				DB:       dbIndex,
			})
			s.dbClients = append(s.dbClients, client)
		}
	} else {
		defaultClient = redis.NewClient(&redis.Options{
			Addr:     config.Config.Host + ":" + strconv.Itoa(config.Config.Port),
			Password: config.Config.Password,
			DB:       config.Config.Db,
		})
		s.dbClients = append(s.dbClients, defaultClient)
	}
	// check if memory usage is available
	_, err := defaultClient.MemoryUsage(ctx, "").Result()
	if err != nil && redis.Nil != err {
		panic(err)
	}
}

func (s *Scanner) setOverallStats(ctx context.Context) {
	if len(s.dbClients) == 0 {
		panic("no db clients, make sure you have setup clients")
	}
	client := s.dbClients[0]
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
	s.wg.Add(1)
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
			if ttlCmd.Err() != nil || memoryUsageCmd.Err() != nil {
				fmt.Printf("error getting ttl or memory usage for key %s: ttlError: %v, memoryUsageError: %v\n", key, ttlCmd.Err(), memoryUsageCmd.Err())
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
	s.setOverallStats(ctx)
	for i := 0; i < config.Config.Concurrency; i++ {
		go s.startScanWorker(ctx)
	}
	for _, client := range s.dbClients {
		for _, keyType := range allKeyTypes {
			var cursor uint64 = 0
			for {
				keys, nextCursor, err := client.ScanType(ctx, cursor, config.Config.ScanPattern, 1000, keyType).Result()
				if err != nil {
					panic(err)
				}
				// Split keys into groups of at most 50 keys each
				for i := 0; i < len(keys); i += 50 {
					end := min(i+50, len(keys))
					group := keys[i:end]

					s.scanJobChan <- scanJob{
						dbClient: client,
						keyType:  keyType,
						keys:     group,
					}
				}
				cursor = nextCursor
				if cursor == 0 {
					break
				}
			}
		}
	}
	close(s.scanJobChan)
	s.wg.Wait()

	defer func() {
		for _, client := range s.dbClients {
			client.Close()
		}
	}()
	defer close(s.finishChan)
}

func (s *Scanner) Finish() chan struct{} {
	return s.finishChan
}
