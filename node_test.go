package main

import (
	"testing"
)

func equal(xs,ys []Interval) bool {
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

func TestLog(t *testing.T) {
	node := NewNode("testLog")

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

func TestPush(t *testing.T) {
	node := NewNode("testPush")

	// Visit some points
	cells := []Point{
		Point{0,0},
		Point{0,1},
		Point{0,2},
		Point{0,3},
		Point{0,4},
	}
	for i, cell := range cells {
		node.Log(cell, i)
	}
	
	node.Push()
	
	// Check initial push
	for i, cell := range cells {
		_, ok := DHT[cell]
		if !ok {
			t.Errorf("Failed at Push(): (cell=%v, time=%v) was missing", cell, i)
		}
	}

}
