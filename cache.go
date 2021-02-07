package shardcache

import (
	"sync"
)

type shard struct {
	mu    sync.Mutex
	items map[uint64]interface{}
}

type ShardCache struct {
	mu     sync.RWMutex
	n      uint64
	shards []shard
}

func New(shards uint64) *ShardCache {
	if shards%2 > 0 {
		panic("shards must be a power of 2")
	}

	return &ShardCache{
		n:      shards,
		shards: make([]shard, shards),
	}
}

func (sl *ShardCache) Get(key uint64) interface{} {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	return sl.shards[key&sl.n].items[key]
}

func (sl *ShardCache) Set(key uint64, value interface{}) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	sl.shards[key&sl.n].mu.Lock()
	defer sl.shards[key&sl.n].mu.Unlock()

	if sl.shards[key&sl.n].items == nil {
		sl.shards[key&sl.n].items = make(map[uint64]interface{})
	}

	sl.shards[key&sl.n].items[key] = value
}

func (sl *ShardCache) Delete(key uint64) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	sl.shards[key&sl.n].mu.Lock()
	defer sl.shards[key&sl.n].mu.Unlock()

	delete(sl.shards[key&sl.n].items, key)
}
