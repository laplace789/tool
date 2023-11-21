package slice

import (
	"reflect"
	"testing"
)

func TestFillBy(t *testing.T) {
	type args[V any] struct {
		slice []V
		fn    func(index int, val V) V
	}
	type testCase[V any] struct {
		name string
		args args[V]
		want []V
	}
	tests := []testCase[int]{
		{
			name: "complete divided",
			args: args[int]{
				slice: []int{-1, -1, -1, -1},
				fn: func(index int, val int) int {
					if index == 0 {
						return 1
					} else {
						return 0
					}
				},
			},
			want: []int{1, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FillBy(tt.args.slice, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FillBy() = %v, want %v", got, tt.want)
			}
		})
	}
}
