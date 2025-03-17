package counter

import (
	"math"
	"sync"
	"time"
)

var (
	overallCounter = OverallCounter{
		TotalMemory: 0,
		UsedMemory:  0,
		TotalKeys:   0,
		ByTTL: []*ByTTL{
			{
				Desc:      "no exp",
				minTTL:    -1,
				maxTTL:    -1,
				TotalKeys: 0,
			},
			{
				Desc:      "expired",
				minTTL:    -2,
				maxTTL:    -2,
				TotalKeys: 0,
			},
			{
				Desc:      "0-1h",
				minTTL:    0,
				maxTTL:    time.Hour,
				TotalKeys: 0,
			},
			{
				Desc:      "1-3h",
				minTTL:    time.Hour + 1,
				maxTTL:    time.Hour * 3,
				TotalKeys: 0,
			},
			{
				Desc:      "3-12h",
				minTTL:    time.Hour*3 + 1,
				maxTTL:    time.Hour * 12,
				TotalKeys: 0,
			},
			{
				Desc:      "12-24h",
				minTTL:    time.Hour*12 + 1,
				maxTTL:    time.Hour * 24,
				TotalKeys: 0,
			},
			{
				Desc:      "1-2day",
				minTTL:    time.Hour*24 + 1,
				maxTTL:    time.Hour * 48,
				TotalKeys: 0,
			},
			{
				Desc:      "3-7day",
				minTTL:    time.Hour*48 + 1,
				maxTTL:    time.Hour * 24 * 7,
				TotalKeys: 0,
			},
			{
				Desc:      ">7day",
				minTTL:    time.Hour*24*7 + 1,
				maxTTL:    math.MaxInt64,
				TotalKeys: 0,
			},
		},
	}
)

func GetOverallCounter() OverallCounter {
	return overallCounter
}

type OverallCounter struct {
	TotalMemory int64
	UsedMemory  int64
	TotalKeys   int64
	ByTTL       []*ByTTL
}

func (c *OverallCounter) FindByTTL(d time.Duration) (ret *ByTTL) {
	for _, b := range c.ByTTL {
		if b.Match(d) {
			return b
		}
	}
	return
}

type ByTTL struct {
	mutex     sync.Mutex
	Desc      string
	minTTL    time.Duration
	maxTTL    time.Duration
	TotalKeys int64
}

func (b *ByTTL) Match(ttl time.Duration) bool {
	return ttl >= b.minTTL && ttl <= b.maxTTL
}

func (b *ByTTL) Add() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.TotalKeys += 1
}
