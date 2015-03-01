package numeric

import (
	"github.com/ready-steady/hotspot"
)

// Config represents the configuration of a particular problem.
type Config struct {
	// The configuration of the HotSpot model.
	hotspot.Config

	// The temperature of the ambience.
	Ambience float64 // in Kelvin
}
