package slice

import (
	"reflect"
	"testing"
)

func TestChunk(t *testing.T) {
	type args[V any] struct {
		slice []V
		size  int
	}
	type testCase[V any] struct {
		name    string
		args    args[V]
		want    [][]V
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name: "complete divided",
			args: args[int]{
				slice: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size:  3,
			},
			want: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			wantErr: false,
		},
		{
			name: "not complete divided",
			args: args[int]{
				slice: []int{1, 2, 3, 4, 5, 6, 7, 8},
				size:  3,
			},
			want: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Chunk(tt.args.slice, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chunk() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Chunk() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkChunk(b *testing.B) {
	var slice = genLageIntSlice()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Chunk(slice, 3)
	}
}

func TestChunkV2(t *testing.T) {
	type args[V any] struct {
		slice []V
		size  int
	}
	type testCase[V any] struct {
		name    string
		args    args[V]
		want    [][]V
		wantErr bool
	}
	tests := []testCase[int]{
		{
			name: "complete divided",
			args: args[int]{
				slice: []int{1, 2, 3, 4, 5, 6, 7, 8, 9},
				size:  3,
			},
			want: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8, 9},
			},
			wantErr: false,
		},
		{
			name: "not complete divided",
			args: args[int]{
				slice: []int{1, 2, 3, 4, 5, 6, 7, 8},
				size:  3,
			},
			want: [][]int{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ChunkV2(tt.args.slice, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ChunkV2() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ChunkV2() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkChunkV2(b *testing.B) {
	var slice = genLageIntSlice()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ChunkV2(slice, 5)
	}
}
