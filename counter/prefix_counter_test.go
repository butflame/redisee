package counter

import (
	"fmt"
	"testing"
)

func Test_countWithPrefix(t *testing.T) {
	type args struct {
		key        string
		memoryUsed int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				key:        "test1",
				memoryUsed: 100,
			},
		},
		{
			name: "test2",
			args: args{
				key:        "test1:test2",
				memoryUsed: 100,
			},
		},
		{
			name: "test3",
			args: args{
				key:        "test1:test2:test3",
				memoryUsed: 100,
			},
		},
		{
			name: "test4",
			args: args{
				key:        "test1:test2:test3:test4",
				memoryUsed: 100,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			countWithPrefix(tt.args.key, tt.args.memoryUsed)
		})
	}
	g := prefixCounter
	fmt.Println(g)
}
