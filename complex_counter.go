package metrics

import (
	"sync/atomic"
)


// ComplexCounter hold an int64 value that can be incremented and decremented.
type ComplexCounter interface {
	Clear()
	Count() int64
	Dec(int64)
	Inc(int64)
}

// GetOrRegisterComplexCounter returns an existing ComplexCounter or constructs and registers
// a new StandardCounter.
func GetOrRegisterComplexCounter(name string, r Registry) ComplexCounter {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewComplexCounter).(ComplexCounter)
}

// NewComplexCounter constructs a new StandardComplexCounter.
func NewComplexCounter() ComplexCounter {
	if UseNilMetrics {
		return NilComplexCounter{}
	}
	return &StandardComplexCounter{0}
}

// NewRegisteNewRegisteredComplexCounter redCounter constructs and registers a new StandardComplexCounter.
func NewRegisteredComplexCounter(name string, r Registry) ComplexCounter {
	c := NewComplexCounter()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// NilCounter is a no-op Counter.
type NilComplexCounter struct{}

// Clear is a no-op.
func (NilComplexCounter) Clear() {}

// Count is a no-op.
func (NilComplexCounter) Count() int64 { return 0 }

// Dec is a no-op.
func (NilComplexCounter) Dec(i int64) {}

// Inc is a no-op.
func (NilComplexCounter) Inc(i int64) {}

// StandardComplexCounter is the standard implementation of a Counter and uses the
// sync/atomic package to manage a single int64 value.
type StandardComplexCounter struct {
	count int64
}

// Clear sets the counter to zero.
func (c *StandardComplexCounter) Clear() {
	atomic.StoreInt64(&c.count, 0)
}

// Count returns the current count.
func (c *StandardComplexCounter) Count() int64 {
	return atomic.LoadInt64(&c.count)
}

// Dec decrements the counter by the given amount.
func (c *StandardComplexCounter) Dec(i int64) {
	atomic.AddInt64(&c.count, -i)
}

// Inc increments the counter by the given amount.
func (c *StandardComplexCounter) Inc(i int64) {
	atomic.AddInt64(&c.count, i)
}
