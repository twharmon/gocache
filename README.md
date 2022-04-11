# Gocache

![](https://github.com/twharmon/gocache/workflows/Test/badge.svg) [![](https://goreportcard.com/badge/github.com/twharmon/gocache)](https://goreportcard.com/report/github.com/twharmon/gocache) [![codecov](https://codecov.io/gh/twharmon/gocache/branch/main/graph/badge.svg?token=K0P59TPRAL)](https://codecov.io/gh/twharmon/gocache)

Thread safe, generic, in-memory cache for Golang with optional TTL settings.

## Documentation
For full documentation see [pkg.go.dev](https://pkg.go.dev/github.com/twharmon/gocache).

## Usage

### Basic
```go
package main

import (
	"github.com/twharmon/gocache"
)

func main() {
    // Create a basic config.
    cfg := gocache.NewConfig().
        WithMaxCapacity(1000).
        WithDefaultEvictionPolicy(time.Second, cocache.NewEvictionPolicy().UpdateOnGet())

    // Create a cache with that config.
    cache := gocache.New[string, int](cfg)

    // Set a value.
    cache.Set("foo", 3)

    // Get a value.
    fmt.Println(cache.Get("foo")) // 3

    // Wait for the value to expire.
    time.Sleep(time.Millisecond * 1001)

    // Get() returns zero value if not found in cache.
    fmt.Println(cache.Get("foo")) // 0

    // Check if cache has a value for a key.
    fmt.Println(cache.Has("foo")) // false
}
```

## Contribute
Make a pull request.