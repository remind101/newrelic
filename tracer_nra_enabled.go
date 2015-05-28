// +build nra_enabled

// The new relic agent sdk currently only support linux (https://docs.newrelic.com/docs/agents/agent-sdk/getting-started/new-relic-agent-sdk)
package nra

import "github.com/remind101/nra/sdk"

// Init initializes the embedded newrelic agent with the given app name and license key.
func Init(app, key string) {
	if _, err := sdk.InitEmbeddedMode(key, app); err != nil {
		panic(err)
	}
}

// NRTxTracer implements the TxTracer interface. It wraps the newrelic package.
type NRTxTracer struct{}

func (t *NRTxTracer) BeginTransaction() (int64, error) {
	return sdk.TransactionBegin()
}
func (t *NRTxTracer) SetTransactionName(txnID int64, name string) error {
	_, err := sdk.TransactionSetName(txnID, name)
	return err
}
func (t *NRTxTracer) BeginGenericSegment(txnID int64, parentID int64, name string) (int64, error) {
	return sdk.SegmentGenericBegin(txnID, parentID, name)
}
func (t *NRTxTracer) BeginDatastoreSegment(txnID int64, parentID int64, table string, operation string, sql string, rollupName string) (int64, error) {
	return sdk.SegmentDatastoreBegin(txnID, parentID, table, operation, sql, rollupName)
}
func (t *NRTxTracer) BeginExternalSegment(txnID int64, parentID int64, host string, name string) (int64, error) {
	return sdk.SegmentExternalBegin(txnID, parentID, host, name)
}
func (t *NRTxTracer) EndSegment(txnID int64, parentID int64) error {
	_, err := sdk.SegmentEnd(txnID, parentID)
	return err
}
func (t *NRTxTracer) SetTransactionRequestURL(txnID int64, url string) error {
	_, err := sdk.TransactionSetRequestURL(txnID, url)
	return err
}
func (t *NRTxTracer) EndTransaction(txnID int64) error {
	_, err := sdk.TransactionEnd(txnID)
	return err
}
