package main

import (
	"testing"
)

func equal(xs,ys []Interval) bool {
	for i, x := range xs {
		if x != ys[i] {
			return false
		}
	}
	return true
}

func TestLog(t *testing.T) {
	node := NewNode("test")

	// Logging interval-points
	visited := []Point{
		Point{0,0}, 
		Point{0,10}, 
		Point{10,10},
		Point{10,10},
		Point{10,10},
		Point{0,10}, 
		Point{0,0}, 
	}
	for i, visit := range visited {
		node.Log(visit, i)
	}
	expected := map[Point][]Interval{
		Point{0,0}: []Interval{Interval{0,0}, Interval{6,6}},
		Point{0,10}: []Interval{Interval{1,1}, Interval{5,5}},
		Point{10,10}: []Interval{Interval{2,4}},
	}
	for p, intervals := range expected {
		if !equal(intervals, node.log[p]) {
			t.Errorf("Failed at Log(): %v != %v", intervals, node.log[p])
		}
	}
}
