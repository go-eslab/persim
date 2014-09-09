package time

import (
	"github.com/go-eslab/persim/system"
)

// List represents a list scheduler.
type List struct {
	plat  *system.Platform
	app   *system.Application
	roots []uint16
}

// NewList creates a new list scheduler for the given platform and application.
func NewList(plat *system.Platform, app *system.Application) *List {
	return &List{
		plat:  plat,
		app:   app,
		roots: app.Roots(),
	}
}

// Compute constructs a schedule according to the given priority vector.
// The length of this vector equals to the number of tasks in the system, and
// smaller values correspond to higher priorities.
func (l *List) Compute(priority []float64) *Schedule {
	cores := l.plat.Cores
	tasks := l.app.Tasks

	cc := uint16(len(cores))
	tc := uint16(len(tasks))

	mapping := make([]uint16, tc)
	order := make([]uint16, tc)
	start := make([]float64, tc)
	finish := make([]float64, tc)

	scheduled := make([]bool, tc)

	ctime := make([]float64, cc)
	ttime := make([]float64, tc)

	var i, j, k, cid, tid, kid, pid uint16
	var ready bool

	size := uint16(len(l.roots))

	// According to the benchmarks, keeping it sorted is not worth it.
	pool := make([]uint16, size, tc)
	copy(pool, l.roots)

	for size > 0 {
		// Find the earliest available core.
		cid = 0
		for i = 1; i < cc; i++ {
			if ctime[i] < ctime[cid] {
				cid = i
			}
		}

		// Find the highest priority task.
		j, tid = 0, pool[0]
		for i = 1; i < size; i++ {
			if priority[pool[i]] < priority[tid] {
				j, tid = i, pool[i]
			}
		}

		// Remove the task from the pool.
		copy(pool[j:], pool[j+1:])
		pool = pool[:size-1]

		mapping[tid] = cid
		order[k] = tid; k++
		if ctime[cid] > ttime[tid] {
			start[tid] = ctime[cid]
		} else {
			start[tid] = ttime[tid]
		}
		finish[tid] = start[tid] + cores[cid].Time[tasks[tid].Type]

		scheduled[tid] = true

		// Update the time when the core is again available.
		ctime[cid] = finish[tid]

		for _, kid = range tasks[tid].Children {
			// Update the time when the child can potentially start executing.
			if ttime[kid] < finish[tid] {
				ttime[kid] = finish[tid]
			}

			// Push the child into the pool if it has become ready for
			// scheduling, that is, if all its parents have been scheduled.
			ready = true

			for _, pid = range tasks[kid].Parents {
				if !scheduled[pid] {
					ready = false
					break
				}
			}

			if !ready {
				continue
			}

			pool = append(pool, kid)
		}

		size = uint16(len(pool))
	}

	return &Schedule{
		Mapping: mapping,
		Order:   order,
		Start:   start,
		Finish:  finish,
	}
}

// Recompute constructs a new schedule based on an old one by adding a delay to
// the execution time of the tasks.
func (l *List) Recompute(s *Schedule, delay []float64) *Schedule {
	tasks := l.app.Tasks

	cc := uint16(len(l.plat.Cores))
	tc := uint16(len(tasks))

	start := make([]float64, tc)
	finish := make([]float64, tc)

	ctime := make([]float64, cc)
	ttime := make([]float64, tc)

	var i, cid, tid, kid uint16

	for ; i < tc; i++ {
		tid = s.Order[i]
		cid = s.Mapping[tid]

		if ctime[cid] > ttime[tid] {
			start[tid] = ctime[cid]
		} else {
			start[tid] = ttime[tid]
		}
		finish[tid] = start[tid] + (s.Finish[tid] - s.Start[tid]) + delay[tid]

		ctime[cid] = finish[tid]

		for _, kid = range tasks[tid].Children {
			if ttime[kid] < finish[tid] {
				ttime[kid] = finish[tid]
			}
		}
	}

	return &Schedule{
		// FIXME: Do not be greedy! Make a copy!
		Mapping: s.Mapping,
		Order:   s.Order,
		Start:   start,
		Finish:  finish,
	}
}
