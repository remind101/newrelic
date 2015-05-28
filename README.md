## NewRelic Go Agent

A small convenience layer that sits on top of [newrelic-go-agent](https://github.com/paulsmith/newrelic-go-agent), to make
it easy to create transactions for NewRelic in Go.

## Caveats

This is alpha software. It has not been tested in a production environment, or any environment for that matter.

## Installing

You'll need to [install the nr_agent_sdk first](https://docs.newrelic.com/docs/agents/agent-sdk/installation-configuration/installing-agent-sdk).

This package will only work on linux platforms. It is also disabled by default. To enable it, use the build flag `newrelic_enabled`:

```
go build -tags newrelic_enabled ./...
```

## Example Usage

``` go
import "github.com/remind101/newrelic"

func main() {
    newrelic.Init("My App", "<new relic license key>")

    // Add to a context.Context
    // https://godoc.org/golang.org/x/net/context
    tx := newrelic.NewTx("/my/transaction/name", nil)
    tx.Start()
    defer tx.End()

    ctx := context.Background()
    ctx = newrelic.WithTx(ctx, tx)

}

// Add a segment to the current transaction if one exists
func FindAllUsers(ctx context.Context, ) ([]User, error) {
    tx, ok := newrelic.FromContext(ctx)
    if ok {
        // Start a datastore segment
        tx.StartDatastore("users", "SELECT", "SELECT * from users WHERE id = 1", "FindAllUsers")
        defer tx.EndSegment()
    }

    // look up users, etc
}

// Add as middleware to your http server
// WARNING: Be sure you understand this https://docs.newrelic.com/docs/apm/other-features/metrics/metric-grouping-issues
// The plan is to build some grouping functionality into this package.
func WithNRA(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        tx := newrelic.NewTx("/my/transaction/name (GET)")
        tx.Start()
        defer tx.End()
        next.ServeHTTP(w, r)
    })
}
```



## Software Credits

The development of this software was made possible using the following components:

[newrelic-go-agent](https://github.com/paulsmith/newrelic-go-agent) by Paul Smith 
Licensed Under: Apache License, Version 2.0
