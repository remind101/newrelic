package main

import (
	"fmt"
	"os"
	"time"

	"github.com/remind101/nra/sdk"
)

func main() {
	debug(sdk.InitEmbeddedMode("NewRelic SDK Test", os.Getenv("NEWRELIC_LICENSE_KEY")))

	for i := 0; i < 10; i++ {
		fmt.Printf("Starting Tx: %d\n", i)

		tx, err := sdk.TransactionBegin()
		debug(tx, err)

		errno, err := sdk.TransactionSetName(tx, "GET /users/{id}")
		debug(errno, err)

		segId, err := sdk.SegmentGenericBegin(tx, 0, "go code")

		time.Sleep(100 * time.Millisecond)

		debug(sdk.SegmentEnd(tx, segId))
		debug(sdk.TransactionEnd(tx))
	}

	time.Sleep(61 * time.Second)

	debug(sdk.RequestShutdown("im done"))
}

func debug(errno interface{}, err error) {
	fmt.Printf("errno: %d, err: %v\n", errno, err)
	if err != nil {
		panic(err)
	}
}
