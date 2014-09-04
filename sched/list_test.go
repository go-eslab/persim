package sched

import (
	"fmt"
	"path"
	"testing"

	"github.com/goesd/persim/system"
	"github.com/goesd/support/assert"
)

const (
	fixturePath = "fixtures"
)

func TestListConstruct(t *testing.T) {
	t.Skip()

	plat, app, _ := system.LoadTGFF(findFixture("002_040"))
	prof := system.NewProfile(plat, app)

	list := NewList(plat, app)
	sched := list.Schedule(prof.Mobility)

	mapping := []uint16{
		0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 0,
		1, 0, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0,
	}

	start := []float64{
		0.0000, 0.0230, 0.0350, 0.0910, 0.0470, 0.0790, 0.4220, 0.0680,
		0.0680, 0.1030, 0.1530, 0.0900, 0.1720, 0.3150, 0.4040, 0.1170,
		0.3360, 0.2040, 0.2590, 0.1380, 0.1530, 0.2720, 0.2190, 0.3360,
		0.1850, 0.3640, 0.3040, 0.2920, 0.2340, 0.2590, 0.2040, 0.3510,
		0.4250, 0.3510, 0.3040, 0.4000, 0.4150, 0.3880, 0.3770, 0.3880,
	}

	duration := []float64{
		0.0230, 0.0120, 0.0120, 0.0120, 0.0210, 0.0110, 0.0180, 0.0230,
		0.0110, 0.0140, 0.0190, 0.0150, 0.0130, 0.0210, 0.0110, 0.0210,
		0.0150, 0.0150, 0.0130, 0.0150, 0.0210, 0.0200, 0.0150, 0.0100,
		0.0190, 0.0130, 0.0110, 0.0120, 0.0250, 0.0150, 0.0180, 0.0190,
		0.0140, 0.0130, 0.0260, 0.0220, 0.0100, 0.0160, 0.0110, 0.0120,
	}

	finish := make([]float64, len(start))
	for i := range start {
		finish[i] = start[i] + duration[i]
	}

	assert.Equal(sched.Mapping, mapping, t)
	assert.Equal(sched.Start, start, t)
	assert.Equal(sched.Finish, finish, t)
}

func findFixture(name string) string {
	return path.Join(fixturePath, fmt.Sprintf("%s.tgff", name))
}
