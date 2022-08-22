package util

import "testing"

func TestWeightedRandomS3(t *testing.T) {
	tests := []struct {
		name        string
		weights     interface{}
		possible    []int
		wantErrText string
	}{
		{
			name:        "int arr",
			weights:     []int{1, 2, 3},
			possible:    []int{0, 1, 2},
			wantErrText: "",
		},
		{
			name:        "int8 arr",
			weights:     []int8{1, 2, 3},
			possible:    []int{0, 1, 2},
			wantErrText: "",
		},
		{
			name:        "int16 arr",
			weights:     []int16{1, 2, 3},
			possible:    []int{0, 1, 2},
			wantErrText: "",
		},
		{
			name:        "int32 arr",
			weights:     []int32{1, 2, 3},
			possible:    []int{0, 1, 2},
			wantErrText: "",
		},
		{
			name:        "int64 arr",
			weights:     []int64{1, 2, 3},
			possible:    []int{0, 1, 2},
			wantErrText: "",
		},
		{
			name:        "float64 arr",
			weights:     []float64{1, 2, 3},
			possible:    []int{0, 1, 2},
			wantErrText: "",
		},
		{
			name:        "float32 arr",
			weights:     []float32{1, 2, 3},
			possible:    []int{0, 1, 2},
			wantErrText: "",
		},
		{
			name:        "String arr",
			weights:     []string{"a", "b", "c"},
			possible:    []int{},
			wantErrText: "no support data type string",
		},
		{
			name:        "String",
			weights:     "test",
			possible:    []int{},
			wantErrText: "input is not slice",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := WeightedRandomS3(tt.weights)
			if err != nil {
				if tt.wantErrText != "" {
					if tt.wantErrText != err.Error() {
						t.Errorf("want err = %s,got err = %s", tt.wantErrText, err.Error())
					}
				} else {
					t.Errorf("unexpected error")
				}
			} else {
				contains := false
				for i := 0; i < len(tt.possible); i++ {
					if got == tt.possible[i] {
						contains = true
						break
					}
				}

				if !contains {
					t.Errorf("got val = %v not in possible = %v", got, tt.possible)
				}
			}

		})
	}
}

func Test_transInterFaceToArr(t *testing.T) {

	tests := []struct {
		name    string
		weights interface{}
		want    []float64
	}{
		{
			name:    "float32",
			weights: []float32{1, 2, 3},
			want:    []float64{1, 2, 3},
		},
		{
			name:    "float64",
			weights: []float64{1, 2, 3},
			want:    []float64{1, 2, 3},
		},
		{
			name:    "int",
			weights: []int{1, 2, 3},
			want:    []float64{1, 2, 3},
		},
		{
			name:    "int8",
			weights: []int8{1, 2, 3},
			want:    []float64{1, 2, 3},
		},
		{
			name:    "int16",
			weights: []int16{1, 2, 3},
			want:    []float64{1, 2, 3},
		},
		{
			name:    "int32",
			weights: []int32{1, 2, 3},
			want:    []float64{1, 2, 3},
		},
		{
			name:    "int64",
			weights: []int64{1, 2, 3},
			want:    []float64{1, 2, 3},
		},
		{
			name:    "string",
			weights: []string{"1", "2", "3"},
			want:    []float64{1, 2, 3},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transInterFaceToArr(tt.weights)
			switch tt.name {
			case "float32", "float64", "int", "int8", "int16", "int32", "int64":
				if err != nil {
					t.Errorf("got err %v", err)
				}

				if len(got) == len(tt.want) {
					for i := 0; i < len(got); i++ {
						if tt.want[i] != got[i] {
							t.Errorf("index %v want val = %v,got val = %v", i, tt.want[i], got[i])
						}
					}
				} else {
					t.Errorf("want len = %v,got len = %v", len(tt.want), len(got))
				}
			default:
				{
					if err.Error() != "no support data type string" {
						t.Errorf("except err = %v got err %v", "no support data type string", err.Error())
					}
				}
			}
		})
	}
}
