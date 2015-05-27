// +build nra_enabled

package sdk

// #cgo LDFLAGS: -L/usr/local/lib -lnewrelic-collector-client -lnewrelic-common -lnewrelic-transaction
// #include "newrelic_collector_client.h"
// #include "newrelic_common.h"
// #include "newrelic_transaction.h"
// #include "stdlib.h"
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"unsafe"
)

// errNoMap is a map of errNo codes to error messages.
var errNoMap = map[int]string{
	-0x10001: "other",
	-0x20001: "disabled",
	-0x30001: "invalid param",
	-0x30002: "invalid id",
	-0x40001: "transaction not started",
	-0x40002: "transaction in progress",
	-0x40003: "transaction not named",
}

// errNo returns an error if the errno is < 0
func errNo(i C.int, context string) (int, error) {
	errno := int(i)
	if errno < 0 {
		errMsg := "unknown"
		if e, ok := errNoMap[errno]; ok {
			errMsg = e
		}
		return errno, errors.New(fmt.Sprintf("newrelic[%s]: %s", context, errMsg))
	}
	return errno, nil
}

func errNoLong(i C.long, context string) (int64, error) {
	_, err := errNo(C.int(i), context)
	return int64(i), err
}

// InitEmbeddedMode registers the message handler with the newrelic embedded message handler
// and calls Init.
func InitEmbeddedMode(license string, appName string) (int, error) {
	C.newrelic_register_message_handler((*[0]byte)(C.newrelic_message_handler))
	return doInit(license, appName, "Go", runtime.Version())
}

/**
 * Start the CollectorClient and the harvester thread that sends application
 * performance data to New Relic once a minute.
 *
 * @param license  New Relic account license key
 * @param app_name  name of instrumented application
 * @param language  name of application programming language
 * @param language_version  application programming language version
 * @return  segment id on success, error code on error, else warning code
 */
func doInit(license string, appName string, language string, languageVersion string) (int, error) {
	clicense := C.CString(license)
	defer C.free(unsafe.Pointer(clicense))

	cappName := C.CString(appName)
	defer C.free(unsafe.Pointer(cappName))

	clang := C.CString("Go")
	defer C.free(unsafe.Pointer(clang))

	clangVersion := C.CString(runtime.Version())
	defer C.free(unsafe.Pointer(clangVersion))

	errno := C.newrelic_init(clicense, cappName, clang, clangVersion)
	return errNo(errno, "init")
}

/**
 * Tell the CollectorClient to shutdown and stop reporting application
 * performance data to New Relic.
 *
 * @reason reason for shutdown request
 * @return  0 on success, error code on error, else warning code
 */
func RequestShutdown(reason string) (int, error) {
	creason := C.CString(reason)
	defer C.free(unsafe.Pointer(creason))
	return errNo(C.newrelic_request_shutdown(creason), "request_shutdown")
}

/*
 * Disable/enable instrumentation. By default, instrumentation is enabled.
 *
 * All Transaction library functions used for instrumentation will immediately
 * return when you disable.
 *
 * @param set_enabled  0 to enable, 1 to disable
 */
func EnableInstrumentation(setEnabled int) {
	C.newrelic_enable_instrumentation(C.int(setEnabled))
}

/*
 * Record a custom metric.
 *
 * @param   name  the name of the metric
 * @param   value   the value of the metric
 * @return  0 on success, else negative warning code or error code
 */
func RecordMetric(name string, value float64) (int, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	return errNo(C.newrelic_record_metric(cname, C.double(value)), "record_metric")
}

/*
 * Record CPU user time in seconds and as a percentage of CPU capacity.
 *
 * @param cpu_user_time_seconds  number of seconds CPU spent processing user-level code
 * @param cpu_usage_percent  CPU user time as a percentage of CPU capacity
 * @return  0 on success, else negative warning code or error code
 */
func RecordCPUUsage(userTimeSecs, cpuUsagePerc float64) (int, error) {
	return errNo(C.newrelic_record_cpu_usage(C.double(userTimeSecs), C.double(cpuUsagePerc)), "record_cpu_usage")
}

/*
 * Record the current amount of memory being used.
 *
 * @param memory_megabytes  amount of memory currently being used
 * @return  0 on success, else negative warning code or error code
 */
func RecordMemoryUsage(memMB float64) (int, error) {
	return errNo(C.newrelic_record_memory_usage(C.double(memMB)), "record_memory_usage")
}

/*
 * Identify the beginning of a transaction. By default, transaction type is set
 * to 'WebTransaction' and transaction category is set to 'Uri'. You can change
 * the transaction type using newrelic_transaction_set_type_other or
 * newrelic_transaction_set_type_web. You can change the transaction category
 * using newrelic_transaction_set_category.
 *
 * @return  transaction id on success, else negative warning code or error code
 */
func TransactionBegin() (int64, error) {
	return errNoLong(C.newrelic_transaction_begin(), "transaction_begin")
}

/*
 * Set the transaction type to 'WebTransaction'. This will automatically change
 * the category to 'Uri'. You can change the transaction category using
 * newrelic_transaction_set_category.
 *
 * @param transaction_id  id of transaction
 * @return  0 on success, else negative warning code or error code
 */
func TransactionSetTypeWeb(id float64) (int, error) {
	return errNo(C.newrelic_transaction_set_type_web(C.long(id)), "transaction_set_type_web")
}

/*
 * Set the transaction type to 'OtherTransaction'. This will automatically
 * change the category to 'Custom'. You can change the transaction category
 * using newrelic_transaction_set_category.
 *
 * @param transaction_id  id of transaction
 * @return  0 on success, else negative warning code or error code
 */
func TransactionSetTypeOther(id float64) (int, error) {
	return errNo(C.newrelic_transaction_set_type_other(C.long(id)), "transaction_set_type_other")
}

/*
 * Set transaction category name, e.g. Uri in WebTransaction/Uri/<txn_name>
 *
 * @param transaction_id  id of transaction
 * @param category  name of the transaction category
 * @return  0 on success, else negative warning code or error code
 */
func TransactionSetCategory(id float64, category string) (int, error) {
	ccategory := C.CString(category)
	defer C.free(unsafe.Pointer(ccategory))
	return errNo(C.newrelic_transaction_set_category(C.long(id), ccategory), "transaction_set_category")
}

/*
 * Identify an error that occurred during the transaction. The first identified
 * error is sent with each transaction.
 *
 * @param transaction_id  id of transaction
 * @param exception_type  type of exception that occurred
 * @param error_message  error message
 * @param stack_trace  stacktrace when error occurred
 * @param stack_frame_delimiter  delimiter to split stack trace into frames
 * @return  0 on success, else negative warning code or error code
 */
func TransactionNoticeError(id float64, exceptionType, errorMessage, stackTrace, stackFrameDelim string) (int, error) {
	cexceptionType := C.CString(exceptionType)
	defer C.free(unsafe.Pointer(cexceptionType))

	cerrorMessage := C.CString(errorMessage)
	defer C.free(unsafe.Pointer(cerrorMessage))

	cstackTrace := C.CString(stackTrace)
	defer C.free(unsafe.Pointer(cstackTrace))

	cstackFrameDelim := C.CString(stackFrameDelim)
	defer C.free(unsafe.Pointer(cstackFrameDelim))

	return errNo(C.newrelic_transaction_notice_error(C.long(id), cexceptionType, cerrorMessage, cstackTrace, cstackFrameDelim), "transaction_notice_error")
}

/*
 * Set a transaction attribute. Up to the first 50 attributes added are sent
 * with each transaction.
 *
 * @param transaction_id  id of transaction
 * @param name  attribute name
 * @param value  attribute value
 * @return  0 on success, else negative warning code or error code
 */
func TransactionAddAttribute(id float64, name, value string) (int, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))

	return errNo(C.newrelic_transaction_add_attribute(C.long(id), cname, cvalue), "transaction_add_attribute")
}

/*
 * Set the name of a transaction.
 *
 * @param transaction_id  id of transaction
 * @param name  transaction name
 * @return  0 on success, else negative warning code or error code
 */
func TransactionSetName(id float64, name string) (int, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return errNo(C.newrelic_transaction_set_name(C.long(id), cname), "transaction_set_name")
}

/*
 * Set the request url of a transaction. The query part of the url is
 * automatically stripped from the url.
 *
 * @param transaction_id  id of transaction
 * @param request_url  request url for a web transaction
 * @return  0 on success, else negative warning code or error code
 */
func TransactionSetRequestURL(id float64, url string) (int, error) {
	curl := C.CString(url)
	defer C.free(unsafe.Pointer(curl))

	return errNo(C.newrelic_transaction_set_request_url(C.long(id), curl), "transaction_set_request_url")
}

/*
 * Set the maximum number of trace segments allowed in a transaction trace. By
 * default, the maximum is set to 2000, which means the first 2000 segments in a
 * transaction will create trace segments if the transaction exceeds the
 * trace threshold (4 x apdex_t).
 *
 * @param transaction_id  id of transaction
 * @param max_trace_segments  maximum number of trace segments
 * @return  0 on success, else negative warning code or error code
 */
func TransactionSetMaxTraceSegments(id float64, max int) (int, error) {
	return errNo(C.newrelic_transaction_set_max_trace_segments(C.long(id), C.int(max)), "transaction_set_max_trace_segments")
}

/*
 * Identify the end of a transaction
 *
 * @param transaction_id  id of transaction
 * @return  0 on success, else negative warning code or error code
 */
func TransactionEnd(id float64) (int, error) {
	return errNo(C.newrelic_transaction_end(C.long(id)), "transaction_end")
}

/*
 * Identify the beginning of a segment that performs a generic operation. This
 * type of segment does not create metrics, but can show up in a transaction
 * trace if a transaction is slow enough.
 *
 * @param transaction_id  id of transaction
 * @param parent_segment_id  id of parent segment
 * @param name  name to represent segment
 * @return  segment id on success, else negative warning code or error code
 */
func SegmentGenericBegin(id, parent float64, name string) (int64, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return errNoLong(C.newrelic_segment_generic_begin(C.long(id), C.long(parent), cname), "segment_generic_begin")
}

/*
 * Identify the beginning of a segment that performs a database operation.
 *
 *
 * SQL Obfuscation
 * ===============
 * If you supply the sql_obfuscator parameter with NULL, the supplied SQL string
 * will go through our basic literal replacement obfuscator that strips the SQL
 * string literals (values between single or double quotes) and numeric
 * sequences, replacing them with the ? character. For example:
 *
 * This SQL:
 * 		SELECT * FROM table WHERE ssn=‘000-00-0000’
 *
 * obfuscates to:
 * 		SELECT * FROM table WHERE ssn=?
 *
 * Because our default obfuscator just replaces literals, there could be
 * cases that it does not handle well. For instance, it will not strip out
 * comments from your SQL string, it will not handle certain database-specific
 * language features, and it could fail for other complex cases.
 *
 * If this level of obfuscation is not sufficient, you can supply your own
 * custom obfuscator via the sql_obfuscator parameter.
 *
 * SQL Trace Rollup
 * ================
 * The agent aggregates similar SQL statements together using the supplied
 * sql_trace_rollup_name.
 *
 * To make the most out of this feature, you should either (1) supply the
 * sql_trace_rollup_name parameter with a name that describes what the SQL is
 * doing, such as "get_user_account" or (2) pass it NULL, in which case
 * it will use the sql obfuscator to generate a name.
 *
 * @param transaction_id  id of transaction
 * @param parent_segment_id  id of parent segment
 * @param table  name of the database table
 * @param operation  name of the sql operation
 * @param sql  the sql string
 * @param sql_trace_rollup_name  the rollup name for the sql trace
 * @param sql_obfuscator  a function pointer that takes sql and obfuscates it
 * @return  segment id on success, else negative warning code or error code
 */
func SegmentDatastoreBegin(id, parent float64, table, operation, sql, sqlTraceRollupName string) (int64, error) {
	ctable := C.CString(table)
	defer C.free(unsafe.Pointer(ctable))

	coperation := C.CString(operation)
	defer C.free(unsafe.Pointer(coperation))

	csql := C.CString(sql)
	defer C.free(unsafe.Pointer(csql))

	csqlTraceRollupName := C.CString(sqlTraceRollupName)
	defer C.free(unsafe.Pointer(csqlTraceRollupName))

	return errNoLong(C.newrelic_segment_datastore_begin(
		C.long(id),
		C.long(parent),
		ctable,
		coperation,
		csql,
		csqlTraceRollupName,
		(*[0]byte)(C.newrelic_basic_literal_replacement_obfuscator),
	), "segment_datastore_begin")
}

/*
 * Identify the beginning of a segment that performs an external service.
 *
 * @param transaction_id  id of transaction
 * @param parent_segment_id  id of parent segment
 * @param host  name of the host of the external call
 * @param name  name of the external transaction
 * @return  segment id on success, else negative warning code or error code
 */
func SegmentExternalBegin(id, parent float64, host, name string) (int64, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	chost := C.CString(host)
	defer C.free(unsafe.Pointer(chost))

	return errNoLong(C.newrelic_segment_external_begin(C.long(id), C.long(parent), chost, cname), "segment_external_begin")
}

/*
 * Identify the end of a segment
 *
 * @param transaction_id  id of transaction
 * @param egment_id  id of the segment to end
 * @return  0 on success, else negative warning code or error code
 */
func SegmentEnd(id, segId float64) (int, error) {
	return errNo(C.newrelic_segment_end(C.long(id), C.long(segId)), "segement_end")
}
