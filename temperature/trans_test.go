package temperature

import (
	"testing"

	"github.com/ready-steady/probability"
	"github.com/ready-steady/probability/uniform"
	"github.com/ready-steady/support/assert"
)

func TestComputeTransient(t *testing.T) {
	temperature, _ := Load(findFixture("002.json"))

	cc := uint(2)
	sc := uint(len(fixtureP)) / cc

	Q := temperature.ComputeTransient(fixtureP, sc)

	assert.EqualWithin(Q, fixtureQ, 1e-12, t)
}

func BenchmarkComputeTransient(b *testing.B) {
	temperature, _ := Load(findFixture("032.json"))

	cc := uint(32)
	sc := uint(1000)

	P := probability.Sample(uniform.New(0, 1), cc*sc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		temperature.ComputeTransient(P, sc)
	}
}
