package main

import (
	"testing"
)

func equal(xs, ys []Interval) bool {
	for i := range xs {
		if xs[i] != ys[i] {
			return false
		}
	}
	return true
}

func TestLocalTraceAdd(t *testing.T) {
	trace := NewLocalTrace()
	codes := []PointCode {
		"loc1",
		"loc2",
		"loc3",
	}
	for t, code := range codes {
		trace.Add(code, t)
	}
	expected := map[PointCode][]Interval {
		"loc1" : []Interval{
			Interval{0,0},
		},
		"loc2" : []Interval{
			Interval{1,1},
		},
		"loc3" : []Interval{
			Interval{2,2},
		},
	}
	for _, code := range codes {
		if !equal(trace.data[code], expected[code]) {
			t.Errorf("expected=%v actual=%v", expected[code], trace.data[code])
		}
	}
}

func TestIntersectIntervalLists(t *testing.T) {
	xs := []Interval{ 
		Interval{0,1},
		Interval{2,2},
		Interval{3,5},
	}
	ys := []Interval{
		Interval{0,3},
		Interval{1,2},
		Interval{2,2},
		Interval{5,7},
	}
	expected := []Interval {
		Interval{0,1},
		Interval{1,1},
		Interval{2,2},
		Interval{3,3},
		Interval{5,5},
	}
	results := Intersect(xs, ys)
	for _, x := range expected {
		if !results[x] {
			t.Errorf("expected=%v actual=%v", expected, results)
		}
	}
}
