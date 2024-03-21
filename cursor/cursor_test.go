package cursor

import (
	"sync"
	"testing"
)

func testAdd(wg *sync.WaitGroup, cursor *Cursor) {
	cursor.Increment()
	wg.Done()
}

func TestCursor_Increment(t *testing.T) {
	tests := []struct {
		name  string
		times int
		want  uint64
	}{
		{name: "test1", times: 1000000, want: 1000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wg sync.WaitGroup
			c := newCursor()
			for i := 0; i < tt.times; i++ {
				wg.Add(1)
				go testAdd(&wg, c)
			}
			wg.Wait()
			got := c.AtomicLoad()
			if got != tt.want {
				t.Errorf("Increment() = %v, want %v", got, tt.want)
			}
		})
	}
}
