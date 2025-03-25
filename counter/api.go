package counter

import "time"

func AddKey(key string, keyType string, ttl time.Duration, memoryUsed int64) {
	countWithTTL(key, ttl)
	countWithType(keyType, memoryUsed)
	countBigKey(key, keyType, memoryUsed)
	// todo
}
