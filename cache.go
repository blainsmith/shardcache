package shardcache

import (
	"sync"
	"time"
)

type shard struct {
	mu    sync.Mutex
	items map[uint64]item
}

type item struct {
	value     interface{}
	expiresAt uint64
}

type ShardLock struct {
	mu     sync.RWMutex
	n      uint64
	exp    uint64
	shards []shard
}

type Options struct {
	Shards     uint64
	Expiration time.Duration
}

func New(opts Options) *ShardLock {
	if opts.Shards%2 > 0 {
		panic("shards must be a power of 2")
	}

	return &ShardLock{
		n:      opts.Shards,
		exp:    uint64(opts.Expiration.Nanoseconds()),
		shards: make([]shard, opts.Shards),
	}
}

func (sl *ShardLock) Get(key uint64) interface{} {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	return sl.shards[key&sl.n].items[key].value
}

func (sl *ShardLock) Set(key uint64, value interface{}) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	sl.shards[key&sl.n].mu.Lock()
	defer sl.shards[key&sl.n].mu.Unlock()

	if sl.shards[key&sl.n].items == nil {
		sl.shards[key&sl.n].items = make(map[uint64]item)
	}

	sl.shards[key&sl.n].items[key] = item{value: value, expiresAt: uint64(time.Now().UnixNano()) + sl.exp}
}

func (sl *ShardLock) Delete(key uint64) {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	sl.shards[key&sl.n].mu.Lock()
	defer sl.shards[key&sl.n].mu.Unlock()

	delete(sl.shards[key&sl.n].items, key)
}
