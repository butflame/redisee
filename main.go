package main

import (
	"context"
	"flag"
	"fmt"
	"redisee/config"
	"strconv"
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

	fmt.Println("Start with setting:")
	fmt.Printf("Host: %s\n", config.Config.Host)
	fmt.Printf("Port: %d\n", config.Config.Port)
	fmt.Printf("Password: %s\n", config.Config.Password)
	dbRepr := strconv.Itoa(config.Config.Db)
	if config.Config.Db == config.ALL_DB {
		dbRepr = "All"
	}
	fmt.Printf("Db: %s\n", dbRepr)
	fmt.Printf("Separator: %s\n", config.Config.Separator)
	fmt.Printf("Scan pattern: %s\n", config.Config.ScanPattern)
	fmt.Printf("Concurrency: %d\n", config.Config.Concurrency)
	fmt.Println()
	fmt.Println("you can run with \"-h\" to see available flags")
	fmt.Println()

	ctx := context.Background()
	instance.Run(ctx)
}
