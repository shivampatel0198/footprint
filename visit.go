package main

import (
	"fmt"
	"strings"
)

// Represents a "visit event"
type Visit struct {
	Id string
	Duration Interval
}

func (v Visit) String() string {
	return fmt.Sprintf("visit[id=%s, duration=%s]", v.Id, v.Duration)
}

// A simple implementation of a VisitLog consisting of a slice of Visits
type VisitLogList struct {
	visits []Visit
}

func (vl VisitLogList) Add(v Visit) VisitLogList {
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

func (vl VisitLogList) Last() Visit {
	return vl.visits[len(vl.visits)-1]
}

func (vl VisitLogList) Size() int {
	return len(vl.visits)
}

func (v VisitLogList) Intersect(u VisitLogList) (out VisitLogList) {
	out = VisitLogList{}
	var xs, ys []Visit
	for i,j := 0,0; i < len(xs) && j < len(ys) ; {
		x := xs[i].Duration
		y := ys[j].Duration
		overlap := x.Intersect(y)
		if !overlap.IsEmpty() {
			// Record ID of second interval
			out.Add(Visit{ys[j].Id, overlap})
		}
		// Get rid of interval with earlier endpoint
		if x.End < y.End { 
			i += 1 
		} else {
			j += 1 
		}
	}
	return
}

func (vl VisitLogList) String() string {
	var s strings.Builder
	s.WriteString("Log:[\n\t")
	s.WriteString(vl.visits[0].String())
	for _, visit := range vl.visits[1:] {
		s.WriteString("\n\t")
		s.WriteString(visit.String())
	}
	s.WriteString("]")
	return s.String()
}

