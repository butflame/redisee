package main

import (
	"context"
	"flag"
	"fmt"
	"redisee/config"
)

func main() {
	showHelp := flag.Bool("h", false, "print usages")

	instance := Redisee{}
	flag.StringVar(&config.Config.Host, "host", "127.0.0.1", "the host of the redis server")
	flag.IntVar(&config.Config.Port, "port", 6379, "the port of the redis server")
	flag.StringVar(&config.Config.Password, "password", "", "the password to access the redis server")
	flag.IntVar(&config.Config.Db, "db", config.ALL_DB, "specify the db to scan, scan all dbs by default")
	flag.StringVar(&config.Config.Separator, "sep", "", "the separator to split the keys, use characters neither alphabet nor numeric by default, you can set multiple separators like \"_-:\"")
	flag.StringVar(&config.Config.ScanPattern, "pattern", "*", "the pattern used to scan keys")
	flag.IntVar(&config.Config.Concurrency, "concurrency", 4, "the concurrency when detecting keys")

	flag.Parse()
	if *showHelp {
		flag.PrintDefaults()
		return
	}
	if len(flag.Args()) == 0 {
		fmt.Println("Start with all default setting, you can run with \"-h\" to see available flags")
	}
	ctx := context.Background()
	instance.Run(ctx)
}
