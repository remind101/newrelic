package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/remind101/newrelic"
	"golang.org/x/net/context"
)

func main() {
	initNewRelic()
	runServer()
}

func runServer() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		ctx := context.Background()
		ctx, t := newrelic.TraceRequest(ctx, "homepage", req)
		defer t.Done()

		someFunc(ctx) // Generic Segment
		apiCall(ctx)  // External Segment
		queryDB(ctx)  // Datastore Segment

		fmt.Fprintln(w, "Hello world!")
	})

	http.ListenAndServe(":8080", mux)
}

func initNewRelic() {
	NRAppName := os.Getenv("NEWRELIC_APP_NAME")
	NRLicense := os.Getenv("NEWRELIC_LICENSE_KEY")

	if NRAppName != "" && NRLicense != "" {
		newrelic.Init(NRAppName, NRLicense)
		newrelic.RecordMetrics(60 * time.Second)
	}
}

func someFunc(ctx context.Context) {
	// Adds a generic segment named `main.someFunc`
	t := newrelic.TraceFunc(ctx)
	defer t.Done()

	// Do stuff
}

func apiCall(ctx context.Context) {
	// Add an external segment
	t := newrelic.TraceExternal(ctx, "http://api.coolapp.com", "coolapp.com")
	defer t.Done()

	// Hit api
}

func queryDB(ctx context.Context) {
	t := newrelic.TraceDatastore(ctx, "users", "SELECT", "select * from users", "select.users")
	defer t.Done()

	// Query DB
}
