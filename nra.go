package nra

import (
	"errors"
	"runtime"

	"github.com/remind101/nra/newrelic-go-agent/newrelic"
)

var ErrTxAlreadyStarted = errors.New("transaction already started")

type Agent interface {
	TxTracer
	Init(license string, appName string, lang string, langVersion string) error
	RequestShutdown(reason string) error
	RecordMetric(name string, val float64) error
}

type TxTracer interface {
	BeginTransaction() int64
	SetTransactionName(txnID int64, name string) error
	BeginGenericSegment(txnID int64, parentID int64, name string) int64
	BeginDatastoreSegment(txnID int64, parentID int64, table string, operation string, sql string, rollupName string) int64
	BeginExternalSegment(txnID int64, parentID int64, host string, name string) int64
	EndSegment(txnID int64, parentID int64) error
	SetTransactionRequestURL(txnID int64, url string) error
	EndTransaction(txnID int64) error
}

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
