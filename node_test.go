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
		if !contains(node.log[p], intervals) {
			t.Errorf("Failed at Log(): %v does not contain %v", node.log[p], intervals)
		}
	}
}

// Check whether vl contains all input intervals
func contains(vl VisitLogList, intervals []Interval) bool {
	for _, v := range vl.visits {
		for _, interval := range intervals {
			if v.Duration == interval {
				return true
			}
		}
	}
	return false
}

func TestPush(t *testing.T) {
	node := NewNode("testPush")

	// Setup
	ps := []Point{
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
		_, ok := DHT[p]
		if !ok {
			t.Errorf("Failed at Push(): (cell=%v, time=%v) was missing", p, i)
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
