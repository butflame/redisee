package main

import (
	"context"
	"redisee/config"
	"redisee/reporter"
	"redisee/scanner"
)

type Redisee struct {
	Config config.ConfigT
}

func (r *Redisee) Run(ctx context.Context) {
	theScanner := scanner.New()
	theScanner.Run(ctx)
	<-theScanner.Finish()

	reporter.Report()
	// overallClient := client.New(redis.Options{
	// 	Addr:     r.Config.Host + ":" + strconv.Itoa(r.Config.Port),
	// 	Password: r.Config.Password,
	// 	DB:       0,
	// })
	// defer closeClient(overallClient)
	// infoMap, err := overallClient.InfoMap(ctx).Result()
	// if err != nil {
	// 	panic(err)
	// }
	// overallCounter := counter.GetOverallCounter()
	// totalMemory, _ := strconv.ParseInt(infoMap["Memory"]["maxmemory"], 10, 64)
	// overallCounter.TotalMemory = totalMemory
	// usedMemory, _ := strconv.ParseInt(infoMap["Memory"]["used_memory"], 10, 64)
	// overallCounter.UsedMemory = usedMemory

	// var dbClients []*redis.Client
	// if r.Config.Db == 0 {
	// 	dbClients = append(dbClients, overallClient)
	// } else if r.Config.Db != AllDb {
	// 	dbClients = append(dbClients, client.New(redis.Options{
	// 		Addr:     r.Config.Host + ":" + strconv.Itoa(r.Config.Port),
	// 		Password: r.Config.Password,
	// 		DB:       r.Config.Db,
	// 	}))
	// }
	// // total keys
	// numberRe := regexp.MustCompile(`\d+`)
	// for dbKey, dbInfo := range infoMap["Keyspace"] {
	// 	numbers := numberRe.FindAllString(dbInfo, -1)
	// 	keysCount, _ := strconv.ParseInt(numbers[0], 10, 64)
	// 	overallCounter.TotalKeys += keysCount

	// 	dbIndex, _ := strconv.Atoi(numberRe.FindAllString(dbKey, 1)[0])

	// 	dbClients = append(dbClients, client.New(redis.Options{
	// 		Addr:     r.Config.Host + ":" + strconv.Itoa(r.Config.Port),
	// 		Password: r.Config.Password,
	// 		DB:       dbIndex,
	// 	}))
	// }

	// for _, dbC := range dbClients {
	// 	keys, _, err := dbC.ScanType(ctx, 0, r.Config.ScanPattern, 1000, "string").Result()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	for _, k := range keys {
	// 		result, err := dbC.TTL(ctx, k).Result()
	// 		if err != nil {
	// 			panic(err)
	// 		}
	// 		byTTL := overallCounter.FindByTTL(result)
	// 		if byTTL != nil {
	// 			byTTL.Add()
	// 		}
	// 	}
	// }
}
