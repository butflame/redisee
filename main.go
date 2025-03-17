package main

import (
	"flag"
	"fmt"
)

func main() {
	showHelp := flag.Bool("h", false, "print usages")

	instance := Redisee{}
	flag.StringVar(&instance.Config.Host, "host", "127.0.0.1", "the host of the redis server")
	flag.IntVar(&instance.Config.Port, "port", 6379, "the port of the redis server")
	flag.StringVar(&instance.Config.Password, "password", "", "the password to access the redis server")
	flag.IntVar(&instance.Config.Db, "db", -1, "specify the db to scan, scan all dbs by default")
	flag.StringVar(&instance.Config.Separator, "sep", "", "the separator to split the keys, use characters neither alphabet nor numeric by default, you can set multiple separators like \"_-:\"")
	flag.StringVar(&instance.Config.ScanPattern, "pattern", "", "the pattern used to scan keys")
	flag.IntVar(&instance.Config.Concurrency, "concurrency", 4, "the concurrency when detecting keys")

	flag.Parse()
	if *showHelp {
		flag.PrintDefaults()
		return
	}
	if len(flag.Args()) == 0 {
		fmt.Println("Start with all default setting, you can run with \"-h\" to see available flags")
	}
	instance.Run()
}
