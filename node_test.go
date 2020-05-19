package main

import (
	"testing"
)

func TestLog(t *testing.T) {
	g := NewGlobalTrace()
	node := NewNode(g)
	locs := []Point {
		Point{0,0}, 
		Point{0,1}, 
		Point{1,0},
		Point{1,1},
		Point{1,1},
		Point{0,1}, 
		Point{0,0}, 
		Point{0,0}, 
	}
	for t, loc := range locs {
		node.Log(loc, t)
	}
	expected := map[Point][]Interval {
		Point{0,0}: []Interval{
			Interval{0,0}, 
			Interval{6,7},
		},
		Point{0,1}: []Interval{
			Interval{1,1}, 
			Interval{5,5},
		},
		Point{1,0}: []Interval{
			Interval{2,2}, 
		},
		Point{1,1}: []Interval{
			Interval{3,4},
		},
	}
	for p, intervals := range expected {
		if !Equal(node.trace.data[Encode(p)], intervals) {
			t.Errorf("expected=%v, actual=%v", intervals, node.trace.data[Encode(p)])
		}
	}
}

func TestPush(t *testing.T) {
	g := NewGlobalTrace()
	node := NewNode(g)

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
		_, ok := g.data[Encode(p)]
		if !ok {
			t.Errorf("(cell=%v, time=%v) was missing", p, i)
		}
	}
}

func TestCheck(t *testing.T) {
	g := NewGlobalTrace()

	// Setup infected node
	infected := NewNode(g)
	xs := []Point{
		Point{0,0},
		Point{0,1},
		Point{0,2},
		Point{0,2},
		Point{0,2},
		Point{0,3},
	}
	for t, x := range xs {
		infected.Log(x,t)
	}
	infected.Push()

	// Setup test node
	node := NewNode(g)
	ys := []Point{
		Point{0,2},
		Point{0,2},
		Point{0,2},
		Point{0,2},
	}
	for t, y := range ys {
		node.Log(y,t)
	}
	
	contacts := node.Check()
	expected :=  []Interval {
		Interval{2,3},
	}
	code := Encode(Point{0,2})
	if !Equal(contacts[code], expected) {
		t.Errorf("expected=%v, actual=%v", expected, contacts[code])
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
