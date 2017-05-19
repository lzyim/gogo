# gogo

gogo is an in-memory key-value cache library for Go.

Any object can be cached, for a given duration or forever, and the cache can be safely used by multiple goroutines.

## Installation
```
go get github.com/soloslee/gogo
```

## Usage
```go
import (
    "fmt"
    "github.com/soloslee/gogo"
)

func main() {
    c := cache.New()

    // Set key "key" to hold the string value "hello", with 10 secondes to live.
    // If key already holds a value, it is overwritten, regardless of its type.
    c.Set("key", "hello", 10)

    // Set the value of the key "mykey" to 42, with no expiration time
    // (the item won't be removed until it is re-set, or removed using
    // c.Del("mykey")
    c.Set("mykey", 42, gogo.NoExpiration)

    // Get the value of key "mykey"
    value, found := c.Get("mykey")
    if found {
        fmt.Println(value)
    }

    // Increments the number stored at key "mykey" by ten
    c.Incr("mykey", 10)

    // Decrements the number stored at key "mykey" by one
    c.Decr("mykey", 1)

    // Returns the number of items in the cache
    number := c.Count()

    // Delete all items from the cache
    c.Flush()
}
```

## License
goim is licensed under [MIT license](http://opensource.org/licenses/MIT).