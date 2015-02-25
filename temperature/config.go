package temperature

import (
	"encoding/json"
	"os"
)

// Config captures the configuration of a particular problem.
type Config struct {
	// The floorplan file of the platform to analyze.
	Floorplan string

	// The options specific to the HotSpot model.
	HotSpot struct {
		// A native configuration file (hotspot.config).
		Config string
		// A line of parameters overwriting the parameters in the above file.
		Params string
	}

	// The sampling interval of temperature analysis. It is the time between
	// two successive samples of power or temperature in power or temperature
	// profiles, respectively. In the formulas given in the general description
	// of the package, it is referred to as Δt.
	TimeStep float64 // in seconds

	// The temperature of the ambience.
	AmbientTemp float64 // in Kelvin
}

// LoadConfig reads the configuration from a JSON file.
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &Config{}

	dec := json.NewDecoder(file)
	if err = dec.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}
