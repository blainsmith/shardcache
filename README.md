# Shard Cache

Warning: Do not use this since it is still evolving.

Shard Cache is meant to remove having to lock the entire cache to write a single item to it. Instead the cache is set up into shared based on the key so there can be an unlimited amount of read locks and a write lock per shard. This lets reads be unblocked for all N-1 shards while a single item in a shard is being written.

```go
// set up a cache with 2 shards, Shards MUST be a power of 2
sc := shardcache.New(shardcache.Options{
    Shards: 2,
})

sc.Set(512, "my-512-value")
sc.Set(101, "my-101-value")

val, ok := sc.Get(101).(string)
if !ok {
    // derp
}

fmt.Println(val) // prints "my-101-value"
```