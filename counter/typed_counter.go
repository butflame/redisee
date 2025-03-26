package counter

import "sync"

var (
	typedCounter = TypedCounter{
		KeyCount:       make(map[string]int),
		KeyMemoryUsage: make(map[string]int64),
	}
	typedMutex sync.Mutex
)

// TypedCounter 结构体用于统计不同类型的 key 数量分布和内存占用分布
type TypedCounter struct {
	// keyCount 记录不同类型的 key 数量
	KeyCount map[string]int
	// keyMemoryUsage 记录不同类型的 key 内存占用
	KeyMemoryUsage map[string]int64
}

func GetTypedCounter() TypedCounter {
	return typedCounter
}

func countWithType(keyType string, memoryUsage int64) {
	typedMutex.Lock()
	defer typedMutex.Unlock()

	// 增加该类型的 key 数量
	if _, ok := typedCounter.KeyCount[keyType]; !ok {
		typedCounter.KeyCount[keyType] = 0
	}
	if _, ok := typedCounter.KeyMemoryUsage[keyType]; !ok {
		typedCounter.KeyMemoryUsage[keyType] = 0
	}
	typedCounter.KeyCount[keyType]++
	// 累加该类型的 key 内存占用
	typedCounter.KeyMemoryUsage[keyType] += memoryUsage
}
