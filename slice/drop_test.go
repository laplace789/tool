package slice

import (
	"math/rand"
	"reflect"
	"testing"
)

func genLageIntSlice() []int {
	res := make([]int, 100000000)
	for i := 0; i < len(res); i++ {
		res[i] = rand.Intn(1000000)
	}
	return res
}

func TestDrop(t *testing.T) {
	type args[V any] struct {
		start int
		n     int
		slice []V
	}
	type testCase[V any] struct {
		name string
		args args[V]
		want []V
	}
	tests := []testCase[int]{
		{
			name: "test1",
			args: args[int]{
				start: 1,
				n:     3,
				slice: []int{1, 2, 3, 4, 5},
			},
			want: []int{1, 5},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Drop(tt.args.start, tt.args.n, tt.args.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Drop() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkDrop(b *testing.B) {
	s := genLageIntSlice()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Drop(1, 3, s)
	}
}
