package reporter

import (
	"fmt"
	"redisee/counter"
	"sort"
)

func Report() {
	reportOverallMemoryUsage()
	reportTopKeys()
	reportKeyTypes()
	reportByPrefix()
}

func reportOverallMemoryUsage() {
	// 获取总体内存使用情况
	overallCounter := counter.GetOverallCounter()
	fmt.Println("### Overall:")

	// 打印总体内存使用情况
	// 计算已使用内存的百分比
	var usedPercentage float64
	if overallCounter.TotalMemory != 0 {
		usedPercentage = float64(overallCounter.UsedMemory) / float64(overallCounter.TotalMemory) * 100
		fmt.Printf("Total Memory: %s, Used Memory: %s (%.2f%%), Total Keys: %d\n", formatMemoryUsage(overallCounter.TotalMemory), formatMemoryUsage(overallCounter.UsedMemory), usedPercentage, overallCounter.TotalKeys)
	} else {
		fmt.Printf("Total Memory: %s, Used Memory: %s, Total Keys: %d\n", formatMemoryUsage(overallCounter.TotalMemory), formatMemoryUsage(overallCounter.UsedMemory), overallCounter.TotalKeys)
	}
	fmt.Println()

	fmt.Println("### Total Keys By TTL:")
	for _, byTTL := range overallCounter.ByTTL {
		fmt.Printf("%s: %d keys in total\n", byTTL.Desc, byTTL.TotalKeys)
	}
	fmt.Println()
}

func reportTopKeys() {
	// 获取大 key 计数器
	bigKeyCounter := counter.GetBigKeyCounter()
	// 打印使用内存最多的前 500 个 key
	fmt.Println("### Top100 Keys by Memory Usage:")
	end := min(100, len(bigKeyCounter.TopKeys))
	for _, keyNode := range bigKeyCounter.GetTopKeys()[:end] {
		fmt.Printf("Key: %s, Type: %s, Memory Usage: %s\n", keyNode.Key, keyNode.KeyType, formatMemoryUsage(keyNode.MemoryUsage))
	}
	fmt.Println()
}

func reportKeyTypes() {
	// 获取 key 类型计数器
	typedCounter := counter.GetTypedCounter()
	// 打印 key 类型统计
	fmt.Println("### Statistics By Key Type:")
	for keyType, count := range typedCounter.KeyCount {
		fmt.Printf("Type: %s, Count: %d, Memory Usage: %s\n", keyType, count, formatMemoryUsage(typedCounter.KeyMemoryUsage[keyType]))
	}
	fmt.Println()
}

func reportByPrefix() {
	// 获取 key 前缀计数器
	prefixCounter := counter.GetPrefixCounter()
	// 打印 key 前缀统计
	queue := []*counter.PrefixNode{&prefixCounter}
	toReportNodes := []*counter.PrefixNode{}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		for _, child := range node.Children {
			child.Prefix = node.Prefix + child.Prefix
			queue = append(queue, child)
		}
		if node.Root || len(node.Children) <= 1 {
			continue
		}
		toReportNodes = append(toReportNodes, node)
	}
	// 先按 key 数量排序
	end := min(100, len(toReportNodes))

	sort.Slice(toReportNodes, func(i, j int) bool {
		return toReportNodes[i].KeyCount > toReportNodes[j].KeyCount
	})
	fmt.Println("### Top100 Key Prefix By Key Count:")
	for _, node := range toReportNodes[:end] {
		fmt.Printf("Prefix: %s, Count: %d, Memory Usage: %s\n", node.Prefix, node.KeyCount, formatMemoryUsage(node.MemoryUsed))
	}
	fmt.Println()

	// 再按内存占用排序
	fmt.Println("### Top100 Key Prefix By Memory Usage:")
	sort.Slice(toReportNodes, func(i, j int) bool {
		return toReportNodes[i].MemoryUsed > toReportNodes[j].MemoryUsed
	})
	for _, node := range toReportNodes[:end] {
		fmt.Printf("Prefix: %s, Count: %d, Memory Usage: %s\n", node.Prefix, node.KeyCount, formatMemoryUsage(node.MemoryUsed))
	}
	fmt.Println()

}

func formatMemoryUsage(memoryUsage int64) string {
	if memoryUsage < 1024 {
		return fmt.Sprintf("%d B", memoryUsage)
	}
	if memoryUsage < 1024*1024 {
		return fmt.Sprintf("%.2f KB", float64(memoryUsage)/1024)
	}
	if memoryUsage < 1024*1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(memoryUsage)/(1024*1024))
	}
	return fmt.Sprintf("%.2f GB", float64(memoryUsage)/(1024*1024*1024))
}
