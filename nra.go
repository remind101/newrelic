package nra

import (
	"errors"
)

var ErrTxAlreadyStarted = errors.New("transaction already started")

type Agent interface {
	TxTracer
	Init(license string, appName string, lang string, langVersion string) error
	RequestShutdown(reason string) error
	RecordMetric(name string, val float64) error
}

type TxTracer interface {
	BeginTransaction() (int64, error)
	SetTransactionName(txnID int64, name string) error
	BeginGenericSegment(txnID int64, parentID int64, name string) (int64, error)
	BeginDatastoreSegment(txnID int64, parentID int64, table string, operation string, sql string, rollupName string) (int64, error)
	BeginExternalSegment(txnID int64, parentID int64, host string, name string) (int64, error)
	EndSegment(txnID int64, parentID int64) error
	SetTransactionRequestURL(txnID int64, url string) error
	EndTransaction(txnID int64) error
}
