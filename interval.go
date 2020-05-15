package main

import (
	"fmt"
)

type Interval struct {
	Start Time
	End Time
}

func (i Interval) String() string {
	if (i.Start == i.End) {
		return fmt.Sprintf("interval[%d]", i.Start)
	}
	return fmt.Sprintf("interval[%d, %d]", i.Start, i.End)
}

func (i Interval) Size() Time {
	return i.End - i.Start
}

func (i Interval) Extend(end Time) Interval {
	return Interval{i.Start, end}
}

func (i Interval) IsEmpty() bool {
	return i.Start == 0 && i.End == 0
}

// Return the intersection of two intervals.  If the intersection is empty,
// return the empty Interval.
func (i Interval) Intersect(j Interval) Interval {
	start, end := max(i.Start, j.Start), min(i.End, j.End)
	if start > end {
		return Interval{}
	}
	return Interval{start, end}
}
