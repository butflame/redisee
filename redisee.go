package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"io"
	"redisee/client"
	"redisee/counter"
	"regexp"
	"strconv"
)

const AllDb = -1

type Redisee struct {
	Config Config
}

type Config struct {
	Host        string
	Port        int
	Password    string
	Db          int
	Separator   string
	ScanPattern string
	Concurrency int
}

func (c *Config) SetDefault() {
	if c.ScanPattern == "" {
		c.ScanPattern = "*"
	}
}

func (r *Redisee) Run(ctx context.Context) {
	overallClient := client.New(redis.Options{
		Addr:     r.Config.Host + ":" + strconv.Itoa(r.Config.Port),
		Password: r.Config.Password,
		DB:       0,
	})
	defer closeClient(overallClient)
	infoMap, err := overallClient.InfoMap(ctx).Result()
	if err != nil {
		panic(err)
	}
	overallCounter := counter.GetOverallCounter()
	totalMemory, _ := strconv.ParseInt(infoMap["Memory"]["maxmemory"], 10, 64)
	overallCounter.TotalMemory = totalMemory
	usedMemory, _ := strconv.ParseInt(infoMap["Memory"]["used_memory"], 10, 64)
	overallCounter.UsedMemory = usedMemory

	var dbClients []*redis.Client
	if r.Config.Db == 0 {
		dbClients = append(dbClients, overallClient)
	} else if r.Config.Db != AllDb {
		dbClients = append(dbClients, client.New(redis.Options{
			Addr:     r.Config.Host + ":" + strconv.Itoa(r.Config.Port),
			Password: r.Config.Password,
			DB:       r.Config.Db,
		}))
	}
	// total keys
	numberRe := regexp.MustCompile(`\d+`)
	for dbKey, dbInfo := range infoMap["Keyspace"] {
		numbers := numberRe.FindAllString(dbInfo, -1)
		keysCount, _ := strconv.ParseInt(numbers[0], 10, 64)
		overallCounter.TotalKeys += keysCount

		dbIndex, _ := strconv.Atoi(numberRe.FindAllString(dbKey, 1)[0])

		dbClients = append(dbClients, client.New(redis.Options{
			Addr:     r.Config.Host + ":" + strconv.Itoa(r.Config.Port),
			Password: r.Config.Password,
			DB:       dbIndex,
		}))
	}

	for _, dbC := range dbClients {
		keys, _, err := dbC.ScanType(ctx, 0, r.Config.ScanPattern, 1000, "string").Result()
		if err != nil {
			panic(err)
		}
		for _, k := range keys {
			result, err := dbC.TTL(ctx, k).Result()
			if err != nil {
				panic(err)
			}
			byTTL := overallCounter.FindByTTL(result)
			if byTTL != nil {
				byTTL.Add()
			}
		}
	}
	// todo
}

func closeClient(closer io.Closer) {
	_ = closer.Close()
}
