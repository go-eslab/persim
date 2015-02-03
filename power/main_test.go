package power

import (
	"path"
	"testing"

	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
	"github.com/ready-steady/support/assert"
)

const (
	fixturePath = "fixtures"
)

func TestCompute(t *testing.T) {
	platform, application, _ := system.Load(findFixture("002_040.tgff"))
	profile := system.NewProfile(platform, application)
	list := time.NewList(platform, application)
	schedule := list.Compute(profile.Mobility)
	power, _ := New(platform, application, 1e-3)

	P := make([]float64, 2*440)
	power.Compute(schedule, P, 440)
	assert.Equal(P, fixturePData, t)

	P = make([]float64, 2*42)
	power.Compute(schedule, P, 42)
	assert.Equal(P, fixturePData[:2*42], t)
}

func BenchmarkCompute(b *testing.B) {
	const (
		Δt = 1e-5
	)

	platform, application, _ := system.Load(findFixture("002_040.tgff"))
	profile := system.NewProfile(platform, application)
	list := time.NewList(platform, application)
	schedule := list.Compute(profile.Mobility)
	power, _ := New(platform, application, Δt)

	sc := uint32(schedule.Span / Δt)
	P := make([]float64, 2*sc)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		power.Compute(schedule, P, sc)
	}
}

func findFixture(name string) string {
	return path.Join(fixturePath, name)
}
