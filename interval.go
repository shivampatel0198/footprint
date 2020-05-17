package main

import (
	"fmt"
)

// Represents an interval of time
// Points are represented as intervals where Lo==Hi
type Interval struct {
	Lo Time
	Hi Time
}

// Return a string representation of an interval
func (i Interval) String() string {
	if (i.Lo == i.Hi) {
		return fmt.Sprintf("interval[%d]", i.Lo)
	}
	return fmt.Sprintf("interval[%d, %d]", i.Lo, i.Hi)
}

func (i Interval) Size() Time {
	return i.Hi - i.Lo
}

// Extends the interval's endpoint.  
// If the extension would shorten the interval, then returns the original interval.
func (i Interval) Extend(hi Time) Interval {
	if hi < i.Hi {
		return i
	}
	return Interval{i.Lo, hi}
}

// Return the intersection of two intervals.  
// If the intersection is empty, ok is false
func (i Interval) Intersect(j Interval) (result Interval, ok bool) {
	lo, hi := Max(i.Lo, j.Lo), Min(i.Hi, j.Hi)
	ok = lo <= hi
	result = Interval{lo, hi}
	return
}

//////////////////////////////
// Implement atree.Interval //
//////////////////////////////

func (i Interval) LowAtDimension(d uint64) int64 {
	return int64(i.Lo)
}

func (i Interval) HighAtDimension(d uint64) int64 {
	return int64(i.Hi)
}

func (i Interval) OverlapsAtDimension(o Interval, d uint64) bool {
	_, ok := i.Intersect(o)
	return ok
}
