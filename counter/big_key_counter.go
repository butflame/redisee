package counter

import (
	"container/heap"
	"sort"
	"sync"
)

type KeyNode struct {
	Key         string
	KeyType     string
	MemoryUsage int64
}

var (
	bigKeyCounter = NewBigKeyCounter(500)
	bigKeyMutex   sync.Mutex
)

func GetBigKeyCounter() BigKeyCounter {
	return *bigKeyCounter
}

type BigKeyCounter struct {
	// 最大堆，用于维护使用内存最大的 500 个 key
	TopKeys []KeyNode
	// 最大容量
	MaxSize int
}

func (h *BigKeyCounter) Len() int { return len(h.TopKeys) }
func (h *BigKeyCounter) Less(i, j int) bool {
	return h.TopKeys[i].MemoryUsage < h.TopKeys[j].MemoryUsage
}
func (h *BigKeyCounter) Swap(i, j int) { h.TopKeys[i], h.TopKeys[j] = h.TopKeys[j], h.TopKeys[i] }
func (h *BigKeyCounter) Push(x interface{}) {
	h.TopKeys = append(h.TopKeys, x.(KeyNode))
}
func (h *BigKeyCounter) Pop() interface{} {
	old := h.TopKeys
	n := len(old)
	x := old[n-1]
	h.TopKeys = old[0 : n-1]
	return x
}

// NewBigKeyCounter 创建一个新的 BigKeyCounter 实例
func NewBigKeyCounter(maxSize int) *BigKeyCounter {
	c := &BigKeyCounter{
		TopKeys: []KeyNode{},
		MaxSize: maxSize,
	}
	heap.Init(c)
	return c
}

// AddKey 向计数器中添加一个新的 key
func (c *BigKeyCounter) AddKey(key string, keyType string, memoryUsage int64) {
	bigKeyMutex.Lock()
	defer bigKeyMutex.Unlock()

	if len(c.TopKeys) < c.MaxSize {
		heap.Push(c, KeyNode{
			Key:         key,
			KeyType:     keyType,
			MemoryUsage: memoryUsage,
		})
	} else if memoryUsage > c.TopKeys[0].MemoryUsage {
		heap.Pop(c)
		heap.Push(c, KeyNode{
			Key:         key,
			KeyType:     keyType,
			MemoryUsage: memoryUsage,
		})
	}
}

func countBigKey(key string, keyType string, memoryUsage int64) {
	bigKeyCounter.AddKey(key, keyType, memoryUsage)
}

// GetTopKeys 获取使用内存最大的 500 个 key
func (c *BigKeyCounter) GetTopKeys() []KeyNode {
	bigKeyMutex.Lock()
	defer bigKeyMutex.Unlock()

	result := make([]KeyNode, len(c.TopKeys))
	copy(result, c.TopKeys)
	sort.Slice(result, func(i, j int) bool {
		return result[i].MemoryUsage > result[j].MemoryUsage
	})
	return result
}
