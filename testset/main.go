package main

import (
	"context"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type insertFunc func(redis.Pipeliner)

const insertBatch = 300

func main() {
	ctx := context.Background()

	c := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   5,
	})

	insertChan := make(chan insertFunc, 10000)
	var wg sync.WaitGroup

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			group := []insertFunc{}
			for f := range insertChan {
				group = append(group, f)
				if len(group) == insertBatch {
					p := c.Pipeline()
					for _, f := range group {
						f(p)
					}
					p.Exec(ctx)
					group = []insertFunc{}
				}
			}
			wg.Done()
		}()
	}

	for i := 0; i < rand.Intn(200000)+20000; i++ {
		randomNum := rand.Intn(9000000000) + 1000000000 // 生成10位随机数
		// 定义可能的前缀段
		var prefixes = []string{
			"ifl:string",
			"ifl:string:ab",
			"ifl:string:ab2",
			"ifl:string:ab3",
			"ifl:string:ab3:qq",
			"ifl:string:ab3:qq2",
			"ifl:string:ab3:qq23",
			"ifl:string:ab3:qq23:xx2",
			"ifl:string:ab3:qq23:xx1",
			"ifl:string:ab3:qq23:xx:op",
			"ifl:string:ab3:qq23:xx:ox2",
		}

		// 随机选择一个前缀
		randomPrefix := prefixes[rand.Intn(len(prefixes))]

		key := randomPrefix + ":" + strconv.Itoa(randomNum)
		// 生成10 - 20位随机字母和数字的字符串
		const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
		randomLength := rand.Intn(11) + 10
		randomStr := make([]byte, randomLength)
		for i := range randomStr {
			randomStr[i] = charset[rand.Intn(len(charset))]
		}
		value := string(randomStr)

		// 生成600 - 3600秒的随机过期时间
		expiration := rand.Intn(60*60*24*7) + 600

		// 写入Redis
		insertChan <- func(p redis.Pipeliner) {
			p.Set(ctx, key, value, time.Duration(expiration)*time.Second)
		}
		// err := c.Set(ctx, key, value, time.Duration(expiration)*time.Second).Err()
		// if err != nil {
		// 	panic(err)
		// }
	}
	// 插入500个zset类型的key
	for i := 0; i < rand.Intn(500000)+5000; i++ {
		randomNum := rand.Intn(9000000000) + 1000000000 // 生成10位随机数
		// 定义十种不同的前缀
		var prefixes = []string{
			"ifl:zset",
			"ifl:zset:sub",
			"ifl:zset:custom",
			"ifl:zset:sub:deep",
			"ifl:zset:sub:special",
			"ifl:zset:sub:complex",
			"ifl:zset:sub:complex:damn",
			"ifl:zset:sub:complex:damn2",
			"ifl:zset:sub:complex:damn3",
		}

		// 随机选择一个前缀
		randomPrefix := prefixes[rand.Intn(len(prefixes))]
		key := randomPrefix + ":" + strconv.Itoa(randomNum)

		// 生成10 - 20个成员，每个成员带有随机分数
		members := make([]redis.Z, 0)
		memberCount := rand.Intn(11) + 10
		for j := 0; j < memberCount; j++ {
			// 生成10 - 20位随机字母和数字的字符串作为成员
			const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
			randomLength := rand.Intn(11) + 10
			randomStr := make([]byte, randomLength)
			for k := range randomStr {
				randomStr[k] = charset[rand.Intn(len(charset))]
			}
			member := string(randomStr)
			// 生成0 - 100的随机分数
			score := float64(rand.Intn(101))
			members = append(members, redis.Z{Score: score, Member: member})
		}

		// 写入Redis ZSet
		insertChan <- func(p redis.Pipeliner) {
			p.ZAdd(ctx, key, members...)
		}
		// err := c.ZAdd(ctx, key, members...).Err()
		// if err != nil {
		// 	panic(err)
		// }

		// 生成600 - 3600秒的随机过期时间
		expiration := rand.Intn(60*60*24*7) + 600
		insertChan <- func(p redis.Pipeliner) {
			p.Expire(ctx, key, time.Duration(expiration)*time.Second)
		}
		// err = c.Expire(ctx, key, time.Duration(expiration)*time.Second).Err()
		// if err != nil {
		// 	panic(err)
		// }
	}

	// 插入800个set类型的key
	for i := 0; i < rand.Intn(80000)+8000; i++ {
		randomNum := rand.Intn(9000000000) + 1000000000 // 生成10位随机数
		// 定义十种不同的前缀
		var prefixes = []string{
			"ifl:set",
			"ifl:set:apple",
			"ifl:set:apple2",
			"ifl:set:apple2:cow",
			"ifl:set:apple2:cow2",
			"ifl:set:apple2:cow2:dog",
			"ifl:set:apple2:cow2:dog2",
		}

		// 随机选择一个前缀
		randomPrefix := prefixes[rand.Intn(len(prefixes))]
		key := randomPrefix + ":" + strconv.Itoa(randomNum)

		// 生成10 - 20个成员
		members := make([]interface{}, 0)
		memberCount := rand.Intn(11) + 10
		for j := 0; j < memberCount; j++ {
			// 生成10 - 20位随机字母和数字的字符串作为成员
			const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
			randomLength := rand.Intn(11) + 10
			randomStr := make([]byte, randomLength)
			for k := range randomStr {
				randomStr[k] = charset[rand.Intn(len(charset))]
			}
			member := string(randomStr)
			members = append(members, member)
		}

		// 写入Redis Set
		insertChan <- func(p redis.Pipeliner) {
			p.SAdd(ctx, key, members...)
		}
		// err := c.SAdd(ctx, key, members...).Err()
		// if err != nil {
		// 	panic(err)
		// }

		// 生成600 - 3600秒的随机过期时间
		expiration := rand.Intn(60*60*24*7) + 600
		insertChan <- func(p redis.Pipeliner) {
			p.Expire(ctx, key, time.Duration(expiration)*time.Second)
		}
		// err = c.Expire(ctx, key, time.Duration(expiration)*time.Second).Err()
		// if err != nil {
		// 	panic(err)
		// }
	}
	// 插入1200个hash类型的key
	for i := 0; i < rand.Intn(60000)+6000; i++ {
		randomNum := rand.Intn(9000000000) + 1000000000 // 生成10位随机数
		// 定义十种不同的前缀
		var prefixes = []string{
			"ifl:hash",
			"ifl:hash:sub",
			"ifl:hash:sub2",
			"ifl:hash:sub:deep",
			"ifl:hash:sub:deep2",
			"ifl:hash:sub:deep:swim",
			"ifl:hash:sub:deep2:swim",
			"ifl:hash:sub:deep2:swim3",
			"ifl:hash:sub:deep:swim2",
			"ifl:hash:sub:deep:swim:fly",
			"ifl:hash:sub:deep:swim:fly2",
			"ifl:hash:sub2:deep:swim:fly2",
		}

		// 随机选择一个前缀
		randomPrefix := prefixes[rand.Intn(len(prefixes))]
		key := randomPrefix + ":" + strconv.Itoa(randomNum)

		// 生成10 - 20个字段和值
		fields := make(map[string]interface{})
		fieldCount := rand.Intn(11) + 10
		for j := 0; j < fieldCount; j++ {
			// 生成10 - 20位随机字母和数字的字符串作为字段
			const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
			randomLength := rand.Intn(11) + 10
			randomFieldStr := make([]byte, randomLength)
			for k := range randomFieldStr {
				randomFieldStr[k] = charset[rand.Intn(len(charset))]
			}
			field := string(randomFieldStr)

			// 生成10 - 20位随机字母和数字的字符串作为值
			randomValueStr := make([]byte, randomLength)
			for k := range randomValueStr {
				randomValueStr[k] = charset[rand.Intn(len(charset))]
			}
			value := string(randomValueStr)

			fields[field] = value
		}

		// 写入Redis Hash
		insertChan <- func(p redis.Pipeliner) {
			p.HSet(ctx, key, fields)
		}
		// err := c.HSet(ctx, key, fields).Err()
		// if err != nil {
		// 	panic(err)
		// }

		// 生成600 - 3600秒的随机过期时间
		expiration := rand.Intn(60*60*24*7) + 600
		insertChan <- func(p redis.Pipeliner) {
			p.Expire(ctx, key, time.Duration(expiration)*time.Second)
		}
		// err = c.Expire(ctx, key, time.Duration(expiration)*time.Second).Err()
		// if err != nil {
		// 	panic(err)
		// }
	}
	// 插入700个list类型的key
	for i := 0; i < rand.Intn(60000)+6000; i++ {
		randomNum := rand.Intn(9000000000) + 1000000000 // 生成10位随机数

		// 定义五种不同的前缀
		var prefixes = []string{
			"ifl:list",
			"ifl:list:sub",
			"ifl:list:sub2",
			"ifl:list:sub2:xx2",
			"ifl:list:sub2:xx",
			"ifl:list:sub2:xx:gg",
			"ifl:list:sub2:qq",
			"ifl:list:sub3",
			"ifl:list:sub3:qq",
			"ifl:list:sub3:xx2",
			"ifl:list:sub3:xx2:gg4",
			"ifl:list:sub3:xx2:gg5",
		}

		// 随机选择一个前缀
		randomPrefix := prefixes[rand.Intn(len(prefixes))]
		key := randomPrefix + ":" + strconv.Itoa(randomNum)

		// 生成10 - 20个成员
		members := make([]interface{}, 0)
		memberCount := rand.Intn(11) + 10
		for j := 0; j < memberCount; j++ {
			// 生成10 - 20位随机字母和数字的字符串作为成员
			const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
			randomLength := rand.Intn(11) + 10
			randomStr := make([]byte, randomLength)
			for k := range randomStr {
				randomStr[k] = charset[rand.Intn(len(charset))]
			}
			member := string(randomStr)
			members = append(members, member)
		}

		// 写入Redis List
		insertChan <- func(p redis.Pipeliner) {
			p.RPush(ctx, key, members...)
		}
		// err := c.RPush(ctx, key, members...).Err()
		// if err != nil {
		// 	panic(err)
		// }

		// 生成600 - 3600秒的随机过期时间
		expiration := rand.Intn(60*60*24*7) + 600
		insertChan <- func(p redis.Pipeliner) {
			p.Expire(ctx, key, time.Duration(expiration)*time.Second)
		}
		// err = c.Expire(ctx, key, time.Duration(expiration)*time.Second).Err()
		// if err != nil {
		// 	panic(err)
		// }
	}
	close(insertChan)
	wg.Wait()
}
