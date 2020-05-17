package main

import (
	
)

// Type aliases
type PointCode string

// Records the location history for a single node.
// The record of visits for a given location is a collection of non-overlapping intervals.
type LocalTrace struct {
	data map[PointCode][]Interval
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
	if ok && t - intervals[len(intervals)-1].Size() == TIME_STEP {
		trace.data[cell] = append(intervals[:len(intervals)-1], intervals[len(intervals)-1].Extend(t))
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

// Records the visit history for all nodes in the network.
// The record of visits for a given location is a collection of potentially overlapping intervals.
type GlobalTrace struct {
	data map[PointCode][]Interval
}

func NewGlobalTrace() *GlobalTrace {
	trace := new(GlobalTrace)
	trace.data = make(map[PointCode][]Interval)
	return trace
}

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
