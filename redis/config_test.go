package rdb

import (
	"reflect"
	"testing"
)

func Test_initRedisConf(t *testing.T) {
	tests := []struct {
		name string
		want *Config
	}{
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initRedisConf(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initRedisConf() = %v, want %v", got, tt.want)
			}
		})
	}
}
