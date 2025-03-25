package counter

import (
	"redisee/config"
	"unicode"
)

var (
	prefixCounter = &PrefixNode{
		Root:     true,
		Children: make(map[string]*PrefixNode),
	}
	givenSeparator map[string]bool = nil
)

func GetPrefixCounter() PrefixNode {
	return *prefixCounter
}

type PrefixNode struct {
	Root       bool
	Prefix     string
	KeyCount   int
	MemoryUsed int64
	Children   map[string]*PrefixNode
}

func countWithPrefix(key string, memoryUsed int64) {
	var visiting *PrefixNode = prefixCounter
	for {
		prefix, suffix := splitKey(key)
		node, ok := visiting.Children[prefix]
		if !ok {
			node = &PrefixNode{
				Root:       false,
				Prefix:     prefix,
				KeyCount:   1,
				MemoryUsed: memoryUsed,
				Children:   make(map[string]*PrefixNode),
			}
			visiting.Children[prefix] = node
			visiting = node
		} else {
			node.KeyCount++
			node.MemoryUsed += memoryUsed
			visiting = node
		}
		if suffix == "" {
			break
		} else {
			key = suffix
		}
	}
}

func splitKey(key string) (prefix string, suffix string) {
	for i, r := range key {
		if config.Config.Separator != "" {
			if givenSeparator == nil {
				givenSeparator = make(map[string]bool)
				for _, s := range config.Config.Separator {
					givenSeparator[string(s)] = true
				}
			}
			if _, ok := givenSeparator[string(r)]; ok {
				prefix = key[:i+1]
				suffix = key[i+1:]
				break
			}
		} else {
			// 未设置分隔符，取非字母且非数字的字符作为分隔符
			if !unicode.IsNumber(r) && !unicode.IsLetter(r) {
				prefix = key[:i+1]
				suffix = key[i+1:]
				break
			}
		}
	}
	if prefix == "" {
		prefix = key
	}
	return
}
