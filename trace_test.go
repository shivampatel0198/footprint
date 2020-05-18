package main

import (
	"testing"
)

func TestLocalTraceAdd(t *testing.T) {
	trace := NewLocalTrace()
	codes := []PointCode {
		"loc1",
		"loc2",
		"loc2",
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
			Interval{1,3},
		},
		"loc3" : []Interval{
			Interval{4,4},
		},
	}
	for _, code := range codes {
		if !Equal(trace.data[code], expected[code]) {
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

func TestLocalTraceIntersectGlobal(t *testing.T) {

	// Setup two local traces u,v
	u, v := NewLocalTrace(), NewLocalTrace()
	xs := []PointCode{
		"loc1",
		"loc1",
		"loc1",
		"loc5",
	}
	ys := []PointCode{
		"loc1",
		"loc1",
		"loc3",
		"loc5",
	}
	for t, x := range xs {
		u.Add(x, t)
	}
	for t, y := range ys {
		v.Add(y, t)
	}

	// Add u's traces to global
	global := NewGlobalTrace()
	u.Iterate(func(c PointCode, xs []Interval) {
		global.Add(c,xs)
	})

	// Intersect global with v
	overlap := v.Intersect(global)
	
	expected := map[PointCode][]Interval {
		"loc1" : []Interval{
			Interval{0,1},
		},
		"loc5" : []Interval{
			Interval{3,3},
		},
	}

	for code, interval := range overlap {
		if !equal(interval, expected[code]) {
			t.Errorf("expected=%v, actual=%v", expected, overlap)
		}
	}

}
