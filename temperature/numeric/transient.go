package numeric

import (
	"github.com/ready-steady/linear/matrix"
)

// Compute calculates the temperature profile corresponding to a power profile.
// The power profile is specified by a function func(time float64, power
// []float64) evaluating the power dissipation at an arbitrary time moment. The
// time moments for which the temperature profile is computed are specified by
// the time array; see the corresponding ODE solver for further details.
//
// http://godoc.org/github.com/ready-steady/numeric/integration/ode#DormandPrince.Compute
func (t *Temperature) Compute(power func(float64, []float64),
	time []float64) ([]float64, []float64, error) {

	nc, nn := t.nc, t.nn

	A, B := t.system.A, t.system.B
	P := make([]float64, nc)

	dSdt := func(t float64, S, dSdt []float64) {
		matrix.Multiply(A, S, dSdt, nn, nn, 1)
		power(t, P)
		for i := uint(0); i < nc; i++ {
			dSdt[i] += B[i] * P[i]
		}
	}

	S, time, err := t.integrator.Compute(dSdt, make([]float64, nn), time)
	if err != nil {
		return nil, nil, err
	}

	ns := uint(len(time))

	Q, Qamb := make([]float64, ns*nc), t.system.Qamb
	for i := uint(0); i < nc; i++ {
		for j := uint(0); j < ns; j++ {
			Q[j*nc+i] = S[j*nn+i] + Qamb
		}
	}

	return Q, time, nil
}
