package util

import (
	"errors"
	"fmt"
	"math/rand"
	"reflect"
)

const (
	intStr     = "int"
	int8Str    = "int8"
	int16Str   = "int16"
	int32Str   = "int32"
	int64Str   = "int64"
	float32Str = "float32"
	float64Str = "float64"
)

func WeightedRandomS3(weights interface{}) (int, error) {
	floatArr, err := transInterFaceToArr(weights)
	if err != nil {
		return -1, err
	}
	return weightedRandomS3(floatArr), nil
}

// algo: https://zliu.org/post/weighted-random/
func weightedRandomS3(weights []float64) int {
	n := len(weights)
	if n == 0 {
		return 0
	}
	cdf := make([]float64, n)
	var sum float64 = 0.0
	for i, w := range weights {
		if i > 0 {
			cdf[i] = cdf[i-1] + w
		} else {
			cdf[i] = w
		}
		sum += w
	}
	r := rand.Float64() * sum
	var l, h int = 0, n - 1
	for l <= h {
		m := l + (h-l)/2
		if r <= cdf[m] {
			if m == 0 || (m > 0 && r > cdf[m-1]) {
				return m
			}
			h = m - 1
		} else {
			l = m + 1
		}
	}
	return -1
}

func transInterFaceToArr(weights interface{}) ([]float64, error) {
	if reflect.TypeOf(weights).Kind() != reflect.Slice {
		return nil, errors.New("input is not slice")
	}

	res := make([]float64, 0)

	s := reflect.ValueOf(weights)
	typName := reflect.TypeOf(weights).Elem().Name()

	//no need to convert
	if typName == float64Str {
		return weights.([]float64), nil
	}

	for i := 0; i < s.Len(); i++ {
		val := s.Index(i).Interface()
		valFloat64, err := convertByTypeName(val, typName)
		if err != nil {
			return nil, err
		}
		res = append(res, valFloat64)
	}

	return res, nil
}

func convertByTypeName(val interface{}, typName string) (float64, error) {
	var valFloat64 float64
	unSupport := false
	switch typName {
	case float32Str:
		//will lose precision if you have many digit
		//https://stackoverflow.com/questions/39642810/why-am-i-losing-precision-while-converting-float32-to-float64
		valFloat64 = float64(val.(float32))
	case intStr:
		valFloat64 = float64(val.(int))
	case int8Str:
		valFloat64 = float64(val.(int8))
	case int16Str:
		valFloat64 = float64(val.(int16))
	case int32Str:
		valFloat64 = float64(val.(int32))
	case int64Str:
		valFloat64 = float64(val.(int64))
	default:
		unSupport = true
	}

	if unSupport {
		return valFloat64, errors.New(fmt.Sprintf("no support data type %v", typName))
	}

	return valFloat64, nil
}
