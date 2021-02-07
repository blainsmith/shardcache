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
	shards []*shard
}

// New returns a ShardCache with the number of shards which must be a power of 2
func New(shards uint64) *ShardCache {
	if shards%2 > 0 {
		panic("shards must be a power of 2")
	}

	sc := ShardCache{
		n:      shards,
		shards: make([]*shard, shards),
	}

	for idx := range sc.shards {
		sc.shards[idx] = &shard{items: make(map[uint64]interface{})}
	}

	return &sc
}

// Get a key from the cache and only use a read lock to access it
func (sc *ShardCache) Get(key uint64) interface{} {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.shards[key&sc.n].items[key]
}

// Set a key in the cache and use a write lock only for the shard that will hold the value leaving all other shards available for reads
func (sc *ShardCache) Set(key uint64, value interface{}) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	sc.shards[key&sc.n].mu.Lock()
	defer sc.shards[key&sc.n].mu.Unlock()

	sc.shards[key&sc.n].items[key] = value
}

// Delete a key in the cache and use a write lock only for the shard that has the value leaving all other shards available for reads
func (sc *ShardCache) Delete(key uint64) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	sc.shards[key&sc.n].mu.Lock()
	defer sc.shards[key&sc.n].mu.Unlock()

	delete(sc.shards[key&sc.n].items, key)
}
