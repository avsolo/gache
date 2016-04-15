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
package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/avsolo/gache/storage"
)

var Store = storage.NewStorage() // Init our Storage

func main() {
    key, i := "number", 1
    Store.Set(key, i, -1) // Init our value
    log.Printf("Lets start with i: %v\n", i)

    // Make and start a simple HTTP Server
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

        i, err := Store.Get(key) // Get value saved before
        if err != nil {
            log.Printf("Unable get var. Key: '%s', Error: %s", key, err.Error())
            return
        }

        err = Store.Update(key, i.(int) + 1, -1) // Now update our value
        if err != nil {
            log.Printf("Unable update. Key: '%s', Error: %#v", key, err.Error())
            return
        }

        newI, err := Store.Get(key) // Check result
        log.Printf("Number saved. Now is: %d", newI.(int))

        w.Write([]byte(fmt.Sprintf("Your number is: %d\n", newI)))
    })

    log.Printf("Starting serve. Open 127.0.0.1:8100 in browser and refresh the page")
    log.Fatal(http.ListenAndServe(":8100", nil))
}
```
