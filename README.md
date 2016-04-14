# Gache

Go in-memory caching libriary. Caching TCP server/client.

##Requirements

No requirements in production mode. For testing and profiling next packages
required:
 - github.com/pkg/profile
 - github.com/stretchr/testify/assert

## Installation

go get github.com/avsolo/gache

# Documentation

## TCP Server

Build for your architecture:

    go build geep.go

Use flags for configuring run mode:

    ./geep [-addr <IP:PORT>] [-log] [-log-dir <path/to/log/dir>]

## From your Go application:

```go
import "github.com/avsolo/gache/storage"
store := storage.NewStorage()

// Set
key := "key"
val := "some value"
ttl := 10 // Time to live in cache
storage.Set(key, val, ttl)

// Get
res, err := storage.Get("key")
if err != nil {
    fmt.Printf("%v", err)
}
```
