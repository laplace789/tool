package hash

import (
	"fmt"
	"strconv"
	"testing"
)

const (
	keySize     = 20
	requestSize = 1000
)

func BenchmarkConsistentHashGet(b *testing.B) {
	ch := NewConsistentHash()
	for i := 0; i < keySize; i++ {
		ch.Add("localhost:" + strconv.Itoa(i))
	}

	for i := 0; i < b.N; i++ {
		ch.Get(i)
	}
}

func Test_ConsistentHashGet(t *testing.T) {
	ch := NewConsistentHash()
	keys := make(map[string]int)
	for i := 0; i < keySize; i++ {
		key := fmt.Sprintf("localhost:" + strconv.Itoa(i))
		ch.Add(key)
		keys[key] = 0
	}
	total := 10000000
	for i := 0; i < 10000000; i++ {
		getKey, ok := ch.Get(i)
		if !ok {
			t.Errorf("no virtual node for index:%d", i)
		}
		_, exists := keys[getKey.(string)]
		if !exists {
			t.Errorf("key is not exists:%s", getKey.(string))
			continue
		}
		keys[getKey.(string)] += 1
	}

	for key, val := range keys {
		t.Logf("key %s count:%d percentage:%f", key, val, val/total)
	}

}
