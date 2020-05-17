package main

import (
	atree "github.com/golang-collections/go-datastructures/augmentedtree"
)

// Type aliases
type PointCode string
type IntervalTree atree.Tree

// Records the location history for a single node.
// The record of visits for a given location is a collection of non-overlapping intervals.
type LocalTrace struct {
	data map[PointCode][]Interval
} 

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
	if !ok {
		trace.data[cell] = make([]Interval,0)
	}
	if ok && t - intervals[len(intervals)-1].Size() == TIME_STEP {
		trace.data[cell] = append(intervals[:len(intervals)-1], intervals[len(intervals)-1].Extend(t))
	} else {
		trace.data[cell] = append(intervals, Interval{t, t})
	}
}

// Records the visit history for all nodes in the network.
// The record of visits for a given location is a collection of potentially overlapping intervals.
type GlobalTrace struct {
	data map[PointCode]IntervalTree
}

func main() {
	
}
