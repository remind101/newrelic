// +build newrelic_enabled

package newrelic

type NRTxReporter struct{}

func (r *NRTxReporter) ReportError(txnID int64, exceptionType, errorMessage, stackTrace, stackFrameDelim string) (int, error) {
	return sdk.TransactionNoticeError(txnID, exceptionType, errorMessage, stackTrace, stackFrameDelim)
}
