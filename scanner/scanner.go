package scanner

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

type Scanner struct {
	mutex      sync.Mutex
	dbClients  []*redis.Client
	finishChan chan struct{}
}

func New() *Scanner {
	scanner := &Scanner{
		finishChan: make(chan struct{}),
	}
	// todo
	return scanner
}

func (s *Scanner) Run(ctx context.Context) {
	// todo
	close(s.finishChan)
}

func (s *Scanner) Finish() chan struct{} {
	return s.finishChan
}
