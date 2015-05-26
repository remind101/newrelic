// +build !nra_enabled

// No op implementation for non linux platforms (new relix agent sdk only support linux right now)
package nra

import (
	"log"
	"runtime"
)

func Init(app, key string) {
	log.Println("Using NoOp NRTxTracer for unspported platform:", runtime.GOOS, runtime.GOARCH)
	return
}

type NRTxTracer struct{}

func (t *NRTxTracer) BeginTransaction() int64 {
	return 0
}
func (t *NRTxTracer) SetTransactionName(txnID int64, name string) error {
	return nil
}
func (t *NRTxTracer) BeginGenericSegment(txnID int64, parentID int64, name string) int64 {
	return 0
}
func (t *NRTxTracer) BeginDatastoreSegment(txnID int64, parentID int64, table string, operation string, sql string, rollupName string) int64 {
	return 0
}
func (t *NRTxTracer) BeginExternalSegment(txnID int64, parentID int64, host string, name string) int64 {
	return 0
}
func (t *NRTxTracer) EndSegment(txnID int64, parentID int64) error {
	return nil
}
func (t *NRTxTracer) SetTransactionRequestURL(txnID int64, url string) error {
	return nil
}
func (t *NRTxTracer) EndTransaction(txnID int64) error {
	return nil
}
