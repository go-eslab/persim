package analytic

import (
	"fmt"
	"testing"

	"github.com/ready-steady/assert"
	"github.com/ready-steady/fixture"
	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/generator"
	"github.com/ready-steady/probability/uniform"
)

func TestFixedNew(t *testing.T) {
	const (
		nc = 2
	)

	temperature := loadFixed(nc)

	assert.Equal(temperature.nc, uint(nc), t)
	assert.Equal(temperature.nn, uint(4*nc+12), t)

	assert.EqualWithin(temperature.D, fixtureD, 1e-14, t)

	assert.EqualWithin(temperature.E, fixtureE, 1e-9, t)
	assert.EqualWithin(temperature.F, fixtureF, 1e-9, t)
}

func TestFixedCompute(t *testing.T) {
	const (
		nc = 2
	)

	temperature := loadFixed(nc)
	Q := temperature.Compute(fixtureP)

	assert.EqualWithin(Q, fixtureQ, 1e-12, t)
}

func BenchmarkFixedCompute002(b *testing.B) {
	const (
		nc = 2
	)

	temperature := loadFixed(nc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(fixtureP)
	}
}

func BenchmarkFixedCompute032(b *testing.B) {
	const (
		nc = 32
		ns = 1000
	)

	temperature := loadFixed(nc)
	P := probability.Sample(uniform.New(0, 1), generator.New(0), nc*ns)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.Compute(P)
	}
}

func loadFixed(nc uint) *Fixed {
	config := &Config{}
	fixture.Load(findFixture(fmt.Sprintf("%03d.json", nc)), config)
	temperature, _ := NewFixed(config)
	return temperature
}
