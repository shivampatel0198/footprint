package main

import (
	"testing"
)

// TODO: test w/ multiple nodes
func TestLog(t *testing.T) {
	node := NewNode("testLog")
	locs := []Point{
		Point{0,0}, 
		Point{0,10}, 
		Point{10,10},
		Point{10,10},
		Point{10,10},
		Point{0,10}, 
		Point{0,0}, 
	}
	for i, loc := range locs {
		node.Log(loc, i)
	}
	expected := map[Point][]Interval {
		Point{0,0}: []Interval{
			Interval{0,0}, 
			Interval{6,6},
		},
		Point{0,10}: []Interval{
			Interval{1,1}, 
			Interval{5,5},
		},
		Point{10,10}: []Interval{
			Interval{2,4},
		},
	}
	for p, intervals := range expected {
		if !Equal(node.trace.data[Encode(p)], intervals) {
			t.Errorf("expected=%v, actual=%v", intervals, node.trace.data[Encode(p)])
		}
	}
}

func TestPush(t *testing.T) {
	node := NewNode("testPush")

	ps := []Point{
		Point{0,0}, 
		Point{0,0}, 
		Point{0,1}, 
		Point{0,2}, 
		Point{0,3}, 
		Point{0,4}, 
	}
	for i, p := range ps {
		node.Log(p, i)
	}
	node.Push()
	
	for i, p := range ps {
		_, ok := bulletin.data[Encode(p)]
		if !ok {
			t.Errorf("(cell=%v, time=%v) was missing", p, i)
		}
	}
}

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
