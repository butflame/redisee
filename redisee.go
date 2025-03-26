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
}
