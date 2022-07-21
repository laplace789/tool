package workerpool

import (
	"sync"
	"testing"
)

func TestNewWorkerPool(t *testing.T) {
	tests := []struct {
		name     string
		priority []int
		input    [][]int
		want     []int
	}{
		{
			name: "test1",
			input: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			priority: []int{1, 2, 3},
			want:     []int{6, 15, 24},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := make([]int, 0)
			var wg sync.WaitGroup
			pool := NewWorkerPool()
			for j := 0; j < 100000; j++ {
				for i := 0; i < len(tt.input); i++ {
					wg.Add(1)
					pool.AddTask(func(params []interface{}) interface{} {
						defer func() {
							wg.Done()
						}()
						sum := 0
						arr := params[0].([]int)
						for k := 0; k < len(arr); k++ {
							sum += arr[k]
						}
						got = append(got, sum)
						return nil
					}, tt.priority[i], tt.input[i])
				}
				for i := 0; i < len(tt.input); i++ {
					wg.Add(1)
					pool.AddTask(func(params []interface{}) interface{} {
						defer func() {
							wg.Done()
						}()
						sum := 0
						arr := params[0].([]int)
						for k := 0; k < len(arr); k++ {
							sum += arr[k]
						}
						got = append(got, sum)
						return nil
					}, tt.priority[i], tt.input[i])
				}
				wg.Wait()
				t.Logf("res = %v", got)
				got = nil
				got = make([]int, 0)
			}
		})
	}
}
