// The new relic agent sdk currently only support linux
package nra

import (
	"runtime"

	"github.com/paulsmith/newrelic-go-agent/newrelic"
)

// Init initializes the embedded newrelic agent with the given app name and license key.
func Init(app, key string) {
	if err := newrelic.Init(key, app, "Go", runtime.Version()); err != nil {
		panic(err)
	}
}

// NRTxTracer implements the TxTracer interface. It wraps the newrelic package.
type NRTxTracer struct{}

func (t *NRTxTracer) BeginTransaction() int64 {
	return newrelic.BeginTransaction()
}
func (t *NRTxTracer) SetTransactionName(txnID int64, name string) error {
	return newrelic.SetTransactionName(txnID, name)
}
func (t *NRTxTracer) BeginGenericSegment(txnID int64, parentID int64, name string) int64 {
	return newrelic.BeginGenericSegment(txnID, parentID, name)
}
func (t *NRTxTracer) BeginDatastoreSegment(txnID int64, parentID int64, table string, operation string, sql string, rollupName string) int64 {
	return newrelic.BeginDatastoreSegment(txnID, parentID, table, operation, sql, rollupName)
}
func (t *NRTxTracer) BeginExternalSegment(txnID int64, parentID int64, host string, name string) int64 {
	return newrelic.BeginExternalSegment(txnID, parentID, host, name)
}
func (t *NRTxTracer) EndSegment(txnID int64, parentID int64) error {
	return newrelic.EndSegment(txnID, parentID)
}
func (t *NRTxTracer) SetTransactionRequestURL(txnID int64, url string) error {
	return newrelic.SetTransactionRequestURL(txnID, url)
}
func (t *NRTxTracer) EndTransaction(txnID int64) error {
	return newrelic.EndTransaction(txnID)
}
