package main

import (
	"fmt"
	"strings"
)

// Represents a "visit event"
type Visit struct {
	VisitorID string
	Duration Interval
}

func (v Visit) String() string {
	return fmt.Sprintf("[%s : %s]", v.VisitorID, v.Duration)
}

// Represents a collection of visits 
type VisitLog interface {

	// Add a Visit to the VisitLog
	Add(v Visit) VisitLog

	// Checks whether a Visit is in the VisitLog
	Contains(v Visit) bool

	// Find visitors that overlap with i, return how long each overlaps
	Intersect(i Interval) map[string]Time

	String() string
}

// A simple implementation of a VisitLog consisting of a list of Visits
// TODO: Write a VisitLog using interval trees
type VisitLogList struct {
	visits []Visit
}

func (vl VisitLogList) Add(v Visit) VisitLog {
	return VisitLogList{append(vl.visits, v)}
}

func (vl VisitLogList) Contains(v Visit) bool {
	for _, visit := range vl.visits {
		if visit == v {
			return true
		}
	}
	return false
}

func (vl VisitLogList) Intersect(i Interval) (out map[string]Time) {
	out = make(map[string]Time)
	for _, v := range vl.visits {
		overlap := i.Intersect(v.Duration)
		if !overlap.IsEmpty() {
			out[v.VisitorID] += overlap.Size()
		}
	}
	return
}

func (vl VisitLogList) String() string {
	var s strings.Builder
	s.WriteString(vl.visits[0].String())
	for _, visit := range vl.visits[1:] {
		s.WriteString("\n")
		s.WriteString(visit.String())
	}
	return s.String()
}

