package nra

import "golang.org/x/net/context"

// Tx represents a transaction.
type Tx struct {
	Tracer TxTracer

	id   int64
	name string
	url  string
	ss   *SegmentStack
}

// NewTx returns a new transaction.
func NewTx(name string) *Tx {
	return &Tx{
		Tracer: &NRTxTracer{},
		name:   name,
		ss:     NewSegmentStack(),
	}
}

// NewRequestTx returns a new transaction with a request url.
func NewRequestTx(name string, url string) *Tx {
	return &Tx{
		Tracer: &NRTxTracer{},
		name:   name,
		url:    url,
		ss:     NewSegmentStack(),
	}
}

// Start starts a transaction, setting the id.
func (t *Tx) Start() error {
	if t.id != 0 {
		return ErrTxAlreadyStarted
	}
	t.id = t.Tracer.BeginTransaction()
	if err := t.Tracer.SetTransactionName(t.id, t.name); err != nil {
		return err
	}
	if t.url != "" {
		return t.Tracer.SetTransactionRequestURL(t.id, t.url)
	}
	return nil
}

// End ends a transaction.
func (t *Tx) End() error {
	for t.ss.Peek() != rootSegment {
		t.EndSegment()
	}
	return t.Tracer.EndTransaction(t.id)
}

// StartGeneric starts a generic segment.
func (t *Tx) StartGeneric(name string) {
	id := t.Tracer.BeginGenericSegment(t.id, t.ss.Peek(), name)
	t.ss.Push(id)
}

// StartDatastore starts a datastore segment.
func (t *Tx) StartDatastore(table, operation, sql, rollupName string) {
	id := t.Tracer.BeginDatastoreSegment(t.id, t.ss.Peek(), table, operation, sql, rollupName)
	t.ss.Push(id)
}

// StartExternal starts an external segment.
func (t *Tx) StartExternal(host, name string) {
	id := t.Tracer.BeginExternalSegment(t.id, t.ss.Peek(), host, name)
	t.ss.Push(id)
}

// EndSegment ends the segment at the top of the stack.
func (t *Tx) EndSegment() error {
	if id, ok := t.ss.Pop(); ok {
		return t.Tracer.EndSegment(t.id, id)
	}
	return nil
}

// WithTx inserts a nra.Tx into the provided context.
func WithTx(ctx context.Context, t *Tx) context.Context {
	return context.WithValue(ctx, txKey, t)
}

// FromContext returns a nra.Tx from the context.
func FromContext(ctx context.Context) (*Tx, bool) {
	t, ok := ctx.Value(txKey).(*Tx)
	return t, ok
}

var key int

const (
	txKey key = iota
)
