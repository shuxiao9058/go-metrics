package metrics

import (
	"math"
	"sync/atomic"
)

// GaugeFloat64s hold a float64 value that can be set arbitrarily.
type ComplexGaugeFloat64 interface {
	Snapshot() ComplexGaugeFloat64
	Update(float64)
	Value() float64
}

// GetOrRegisterGaugeFloat64 returns an existing GaugeFloat64 or constructs and registers a
// new StandardGaugeFloat64.
func GetOrRegisterComplexGaugeFloat64(name string, r Registry) ComplexGaugeFloat64 {
	if nil == r {
		r = DefaultRegistry
	}
	return r.GetOrRegister(name, NewComplexGaugeFloat64()).(ComplexGaugeFloat64)
}

// NewGaugeFloat64 constructs a new StandardGaugeFloat64.
func NewComplexGaugeFloat64() ComplexGaugeFloat64 {
	if UseNilMetrics {
		return NilComplexGaugeFloat64{}
	}
	return &StandardComplexGaugeFloat64{
		value: 0.0,
	}
}

// NewRegisteredGaugeFloat64 constructs and registers a new StandardGaugeFloat64.
func NewRegisteredComplexGaugeFloat64(name string, r Registry) ComplexGaugeFloat64 {
	c := NewComplexGaugeFloat64()
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// NewFunctionalGauge constructs a new FunctionalGauge.
func NewFunctionalComplexGaugeFloat64(f func() float64) ComplexGaugeFloat64 {
	if UseNilMetrics {
		return NilComplexGaugeFloat64{}
	}
	return &FunctionalComplexGaugeFloat64{value: f}
}

// NewRegisteredFunctionalGauge constructs and registers a new StandardGauge.
func NewRegisteredFunctionalComplexGaugeFloat64(name string, r Registry, f func() float64) ComplexGaugeFloat64 {
	c := NewFunctionalComplexGaugeFloat64(f)
	if nil == r {
		r = DefaultRegistry
	}
	r.Register(name, c)
	return c
}

// GaugeFloat64Snapshot is a read-only copy of another GaugeFloat64.
type ComplexGaugeFloat64Snapshot float64

// Snapshot returns the snapshot.
func (g ComplexGaugeFloat64Snapshot) Snapshot() ComplexGaugeFloat64 { return g }

// Update panics.
func (ComplexGaugeFloat64Snapshot) Update(float64) {
	panic("Update called on a ComplexGaugeFloat64Snapshot")
}

// Value returns the value at the time the snapshot was taken.
func (g ComplexGaugeFloat64Snapshot) Value() float64 { return float64(g) }

// NilGauge is a no-op Gauge.
type NilComplexGaugeFloat64 struct{}

// Snapshot is a no-op.
func (NilComplexGaugeFloat64) Snapshot() ComplexGaugeFloat64 { return NilComplexGaugeFloat64{} }

// Update is a no-op.
func (NilComplexGaugeFloat64) Update(v float64) {}

// Value is a no-op.
func (NilComplexGaugeFloat64) Value() float64 { return 0.0 }

// StandardGaugeFloat64 is the standard implementation of a GaugeFloat64 and uses
// sync.Mutex to manage a single float64 value.
type StandardComplexGaugeFloat64 struct {
	value uint64
}

// Snapshot returns a read-only copy of the gauge.
func (g *StandardComplexGaugeFloat64) Snapshot() ComplexGaugeFloat64 {
	return ComplexGaugeFloat64Snapshot(g.Value())
}

// Update updates the gauge's value.
func (g *StandardComplexGaugeFloat64) Update(v float64) {
	atomic.StoreUint64(&g.value, math.Float64bits(v))
}

// Value returns the gauge's current value.
func (g *StandardComplexGaugeFloat64) Value() float64 {
	return math.Float64frombits(atomic.LoadUint64(&g.value))
}

// FunctionalGaugeFloat64 returns value from given function
type FunctionalComplexGaugeFloat64 struct {
	value func() float64
}

// Value returns the gauge's current value.
func (g FunctionalComplexGaugeFloat64) Value() float64 {
	return g.value()
}

// Snapshot returns the snapshot.
func (g FunctionalComplexGaugeFloat64) Snapshot() ComplexGaugeFloat64 { return ComplexGaugeFloat64Snapshot(g.Value()) }

// Update panics.
func (FunctionalComplexGaugeFloat64) Update(float64) {
	panic("Update called on a FunctionalComplexGaugeFloat64")
}
