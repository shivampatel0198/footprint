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
	if i.Lo == i.Hi {
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

// Uniquely encode Hi, Lo as a single natural number
// Note: Not perfect b/c of finite encoding, but will work for now
func cantor_pairing(a, b int) uint64 {
	x, y := float64(a), float64(b)
	return uint64(0.5*(x+y)*(x+y+1) + y)
}

// Note: Two intervals are equivalent if their endpoints are equivalent
func (i Interval) ID() uint64 {
	return cantor_pairing(i.Lo, i.Hi)
}

func Equal(xs, ys []Interval) bool {
	if len(xs) != len(ys) {
		return false
	}
	for i, x := range xs {
		if x != ys[i] {
			return false
		}
	}
	return true
}
