// Package provides an Runge–Kutta integrator for systems of ordinary
// differential equations modeling temperature of electronic systems.
package numeric

import (
	"errors"

	"github.com/ready-steady/hotspot"
	"github.com/ready-steady/numeric/integration/ode"
)

// Temperature represents an integrator.
type Temperature struct {
	Cores uint
	Nodes uint

	system system

	integrator *ode.DormandPrince
}

// New returns an integrator set up according to the given configuration.
func New(c *Config) (*Temperature, error) {
	if c.TimeStep <= 0 {
		return nil, errors.New("the time step should be positive")
	}

	model := hotspot.New((*hotspot.Config)(&c.Config))
	cc, nc := model.Cores, model.Nodes

	// Reusing model.G to store A and model.C to store B.
	A := model.G
	B := model.C
	for i := uint(0); i < nc; i++ {
		B[i] = 1 / model.C[i]
	}
	for i := uint(0); i < nc; i++ {
		for j := uint(0); j < nc; j++ {
			A[j*nc+i] = -B[i] * A[j*nc+i]
		}
	}

	integrator, err := ode.NewDormandPrince(&ode.Config{
		MaxStep:  0,
		TryStep:  0,
		AbsError: 1e-3,
		RelError: 1e-3,
	})
	if err != nil {
		return nil, err
	}

	temperature := &Temperature{
		Cores: cc,
		Nodes: nc,

		system: system{
			A: A,
			B: B,

			Δt:   c.TimeStep,
			Qamb: c.Ambience,
		},

		integrator: integrator,
	}

	return temperature, nil
}
