package rdb

import (
	"context"
	"testing"
	"time"
)

func TestRedisClient_Get(t *testing.T) {

	tests := []struct {
		name    string
		key     string
		val     int
		wantVal string
	}{
		{
			name:    "test1",
			key:     "a",
			val:     1,
			wantVal: "1",
		},
		{
			name:    "test1",
			key:     "b",
			val:     2,
			wantVal: "2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := InitRedisConf("../conf")
			rdb := NewRedisClient(cfg)
			ctx := context.Background()
			errSet := rdb.Set(ctx, tt.key, tt.val, time.Minute*50)
			if errSet != nil {
				t.Errorf("rdb set key err %v", errSet)
			}
			got, err := rdb.Get(ctx, tt.key)
			if err != nil {
				t.Errorf("err = %v", err)
			}
			if got != tt.wantVal {
				t.Errorf("want val = %v,got val = %v", tt.wantVal, got)
			}
		})
	}
}
