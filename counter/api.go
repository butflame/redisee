package counter

import "time"

func SetTotalMemory(memory int64) {
	overallMutex.Lock()
	defer overallMutex.Unlock()

	overallCounter.TotalMemory = memory
}

func SetUsedMemory(memory int64) {
	overallMutex.Lock()
	defer overallMutex.Unlock()

	overallCounter.UsedMemory = memory
}

func AddKey(key string, keyType string, ttl time.Duration, memoryUsed int64) {
	countWithTTL(key, ttl)
	countWithType(keyType, memoryUsed)
	countBigKey(key, keyType, memoryUsed)
	countWithPrefix(key, memoryUsed)
}
