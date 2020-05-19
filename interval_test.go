package main

import (
	"testing"
)

func TestIntervalSize(t *testing.T) {
	tests := []Interval{
		Interval{0, 0},
		Interval{0, 1},
		Interval{1, 2},
	}
	expected := []Time{
		0, 1, 1,
	}
	for i, interval := range tests {
		if interval.Size() != expected[i] {
			t.Errorf("%v is not of size %v", interval, expected[i])
		}
	}
}

func TestIntervalExtend(t *testing.T) {
	tests := []Interval{
		Interval{0, 0},
		Interval{5, 10},
		Interval{0, 10},
		Interval{5, 10},
	}
	his := []Time{
		10, 20, 5, 5,
	}
	expected := []Interval{
		Interval{0, 10},
		Interval{5, 20},
		Interval{0, 10},
		Interval{5, 10},
	}
	for i, interval := range tests {
		result := interval.Extend(his[i])
		if result != expected[i] {
			t.Errorf("expected=%v, actual=%v", expected[i], result)
		}
	}
}

type pair struct {
	a Interval
	b Interval
}

func TestIntervalIntersect(t *testing.T) {
	tests := []pair{
		pair{
			Interval{1, 4},
			Interval{2, 3},
		},
		pair{
			Interval{1, 3},
			Interval{2, 4},
		},
		pair{
			Interval{1, 2},
			Interval{2, 2},
		},
		pair{
			Interval{1, 2},
			Interval{2, 3},
		},
		pair{
			Interval{1, 2},
			Interval{3, 4},
		},
	}
	oks := []bool{
		true, true, true, true, false,
	}
	expected := []Interval{
		Interval{2, 3},
		Interval{2, 3},
		Interval{2, 2},
		Interval{2, 2},
		Interval{},
	}

	for i, pair := range tests {
		a, b := pair.a, pair.b
		result, ok := a.Intersect(b)
		if ok != oks[i] {
			t.Errorf("%v,%v intersect? expected=%v, actual=%v", a, b, oks[i], ok)
		}
		if ok && result != expected[i] {
			t.Errorf("Intersect(%v, %v)? expected=%v, actual=%v", a, b, expected[i], result)
		}
	}
}
