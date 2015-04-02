// Package power provides algorithms for simulating the power dissipation of
// concurrent applications running on multiprocessor platforms.
package power

import (
	"github.com/ready-steady/simulation/system"
	"github.com/ready-steady/simulation/time"
)

// Power represents a power simulator configured for a particular system.
type Power struct {
	platform    *system.Platform
	application *system.Application
}

// New returns a power distributor for a platform and an application.
func New(platform *system.Platform, application *system.Application) *Power {
	return &Power{
		platform:    platform,
		application: application,
	}
}

// Sample computes a power profile of a schedule with a sampling interval Δt.
// The ns parameter controls the number of samples that the output matrix will
// contain; short schedules are extended while long ones are truncated.
func (p *Power) Sample(schedule *time.Schedule, Δt float64, ns uint) []float64 {
	cores, tasks := p.platform.Cores, p.application.Tasks
	nc, nt := uint(len(cores)), uint(len(tasks))

	P := make([]float64, nc*ns)

	if count := uint(schedule.Span / Δt); count < ns {
		ns = count
	}

	for i := uint(0); i < nt; i++ {
		j := schedule.Mapping[i]
		p := cores[j].Power[tasks[i].Type]

		s := uint(schedule.Start[i] / Δt)
		f := uint(schedule.Finish[i] / Δt)
		if f > ns {
			f = ns
		}

		for ; s < f; s++ {
			P[s*nc+j] = p
		}
	}

	return P
}

// Progress takes a schedule and returns a function func(time float64, power
// []float64) that computes the power dissipation at an arbitrary time moment
// according to the schedule.
func (p *Power) Progress(schedule *time.Schedule) func(float64, []float64) {
	cores, tasks := p.platform.Cores, p.application.Tasks
	nc, nt := uint(len(cores)), uint(len(tasks))

	mapping := make([][]uint, nc)
	for i := uint(0); i < nc; i++ {
		mapping[i] = make([]uint, 0, nt)
		for j := uint(0); j < nt; j++ {
			if i == schedule.Mapping[j] {
				mapping[i] = append(mapping[i], j)
			}
		}
	}

	start, finish := schedule.Start, schedule.Finish

	return func(time float64, power []float64) {
		for i := uint(0); i < nc; i++ {
			power[i] = 0
			for _, j := range mapping[i] {
				if start[j] <= time && time <= finish[j] {
					power[i] = cores[i].Power[tasks[j].Type]
					break
				}
			}
		}
	}
}
