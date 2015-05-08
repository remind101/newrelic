package nra

import (
	"runtime"

	"github.com/remind101/nra/newrelic-go-agent/newrelic"
)

type Agent struct{}

// New returns a newrelic agent initialized with the given app name and license key.
func New(app, key string) *Agent {
	if err := newrelic.Init(key, app, "Go", runtime.Version()); err != nil {
		panic(err)
	}
	return &Agent{}
}
