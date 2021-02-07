package shardcache

import (
	"sync"
)

type shard struct {
	mu    sync.Mutex
	items map[uint64]interface{}
}

// ShardCache will keep track and shard keys based on the number of shards it holds around mutexes so this is safe for concurrent use
type ShardCache struct {
	mu     sync.RWMutex
	n      uint64
	shards []shard
}

// New returns a ShardCache with the number of shards which must be a power of 2
func New(shards uint64) *ShardCache {
	if shards%2 > 0 {
		panic("shards must be a power of 2")
	}

	return &ShardCache{
		n:      shards,
		shards: make([]shard, shards),
	}
}

// Get a key from the cache and only use a read lock to access it
func (sl *ShardCache) Get(key uint64) interface{} {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	return sl.shards[key&sl.n].items[key]
}

// Set a key in the cache and use a write lock only for the shard that will hold the value leaving all other shards available for reads
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

// Delete a key in the cache and use a write lock only for the shard that has the value leaving all other shards available for reads
func (sl *ShardCache) Delete(key uint64) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	sl.shards[key&sl.n].mu.Lock()
	defer sl.shards[key&sl.n].mu.Unlock()

	delete(sl.shards[key&sl.n].items, key)
}
