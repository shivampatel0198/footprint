package main

import "strings"

// Type aliases
type PointCode string

// Records the location history for a single node.
// The record of visits for a given location is a collection of non-overlapping intervals.
type LocalTrace struct {
	data map[PointCode][]Interval
} 

func (trace *LocalTrace) String() string {
	return String(trace.data)
}

// Create a new LocalTrace
func NewLocalTrace() *LocalTrace {
	trace := new(LocalTrace)
	trace.data = make(map[PointCode][]Interval)
	return trace
}

// Records a visit to a cell at time t
func (trace *LocalTrace) Add(cell PointCode, t Time) {
	intervals, ok := trace.data[cell]

	// If the last recorded location was in the same place, extend the last interval.
	// Otherwise, append a point interval to the trace.
	if ok && t - intervals[len(intervals)-1].Hi == TIME_STEP {
		trace.data[cell] = append(
			intervals[:len(intervals)-1], 
			intervals[len(intervals)-1].Extend(t),
		)
	} else {
		trace.data[cell] = append(intervals, Interval{t, t})
	}
}

// Apply f on each of the cell-intervals pairs in the LocalTrace
func (trace *LocalTrace) Iterate(f func(cell PointCode, intervals []Interval)) {
	for pc, interval := range trace.data {
		f(pc, interval)
	}
}

// Find all of the intersections between a global trace and a local trace
func (local *LocalTrace) Intersect(global *GlobalTrace) (overlap map[PointCode][]Interval) {
	overlap = make(map[PointCode][]Interval)
	local.Iterate(func(c PointCode, xs []Interval) {
		ys, ok := global.Read(c)
		if ok {
			set := Intersect(xs, ys)
			for x := range set {
				overlap[c] = append(overlap[c], x)
			}
		}
	})
	return
}

// Records the visit history for all nodes in the network.
// The record of visits for a given location is a collection of potentially overlapping intervals.
// TODO: thread safety
type GlobalTrace struct {
	data map[PointCode][]Interval
}

func NewGlobalTrace() *GlobalTrace {
	trace := new(GlobalTrace)
	trace.data = make(map[PointCode][]Interval)
	return trace
}

func (g *GlobalTrace) String() string {
	return String(g.data)
}

// NOTE: Rep exposure via mutability
func (g *GlobalTrace) Read(cell PointCode) (xs []Interval, ok bool) {
	xs, ok = g.data[cell]
	return
}

// Add a collection of intervals to a cell
func (trace *GlobalTrace) Add(cell PointCode, intervals []Interval) {
	for _, interval := range intervals {
		trace.data[cell] = append(trace.data[cell], interval)
	}
}

/* Utilities */

// Returns the set of intersections between two lists of intervals
func Intersect(xs,ys []Interval) (overlaps map[Interval]bool) {
	overlaps = make(map[Interval]bool)
	for _, x := range xs {
		for _, y := range ys {
			overlap, ok := x.Intersect(y)
			if ok {
				overlaps[overlap] = true
			}
		}
	}
	return
}

func String(data map[PointCode][]Interval) string {
	var b strings.Builder
	b.WriteString("Trace:")
	for c, intervals := range data {
		b.WriteString("\n  ")
		b.WriteString(string(c))
		b.WriteString(": ")
		for _, interval := range intervals {
			b.WriteString("\n    ")
			b.WriteString(interval.String())
		}
	}
	return b.String()
}
