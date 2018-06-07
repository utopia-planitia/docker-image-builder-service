## TTLCache - an in-memory cache with expiration

TTLCache is a minimal wrapper over a map in golang, entries of which are

1. Thread-safe
2. Expiring after a certain time
3. Extending expiration on `Get`s

[![Build Status](https://travis-ci.org/damoon/ttlcache.svg)](https://travis-ci.org/damoon/ttlcache)

#### Usage
```go
import (
  "time"
  "github.com/damoon/ttlcache"
)

func main () {
  cache := ttlcache.NewCache(time.Second, "")
  cache.Set("key", "value")
  value, exists := cache.Get("key")
  count := cache.Count()

  intCache := ttlcache.NewCache(time.Second, 0)
  intCache.Set(1, 2)
  value, exists := intCache.Get(1)
  intCache := cache.Count()
}
```