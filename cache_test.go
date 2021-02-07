package shardcache_test

import (
	"testing"

	"blainsmith.com/go/shardcache"
)

func TestShardCache(t *testing.T) {
	sc := shardcache.New(64)

	sc.Set(0, "zero")
	sc.Set(1, 1)
	sc.Set(2, []byte{0x02})

	if val, ok := sc.Get(0).(string); !ok || val != "zero" {
		t.Errorf("could not get key 0 as 'zero', %v, %v", ok, val)
	}

	if val, ok := sc.Get(1).(int); !ok || val != 1 {
		t.Errorf("could not get key 1 as 1, %v, %v", ok, val)
	}

	if val, ok := sc.Get(2).([]byte); !ok || val[0] != 0x02 {
		t.Errorf("could not get key 2 as 0x02, %v, %v", ok, val)
	}
}

func BenchmarkShardCacheGet(b *testing.B) {
	sc := shardcache.New(64)

	sc.Set(0, nil)
	sc.Set(1, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.Get(1)
	}
}

func BenchmarkShardCacheSet(b *testing.B) {
	sc := shardcache.New(64)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sc.Set(1, "test")
	}
}
