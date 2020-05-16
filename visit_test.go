package main

import (
	"testing"
)

func TestContains(t *testing.T) {
	v := Visit{"test", Interval{0,2}}
	vl := []Visit{v}
	vll := VisitLogList{vl}
	if !vll.Contains(v) {
		t.Errorf("Contains() test failed: %v not in VisitLogList", v)
	}
}

func TestIntersect(t *testing.T) {
	u := VisitLogList{
		[]Visit{
			Visit{"jesus", Interval{0, 32}},
			Visit{"shivam@williams", Interval{2016, 2020}},
		},
	}
	v := VisitLogList{
		[]Visit{
			Visit{"shivam", Interval{1998, 2020}},
			Visit{"covid", Interval{2019, 2020}},
		},
	}
	expected := VisitLogList{
		[]Visit{
			Visit{"shivam", Interval{2016, 2020}},
			Visit{"covid", Interval{2019, 2020}},
		},
	}
	if !equals(u.Intersect(v), expected) {
		t.Error("Intersect() test failed")
	}
}

func TestExtend(t *testing.T) {
	vl := VisitLogList{}
	vl.Add(Visit{"a", Interval{0,5}})
	vl.Add(Visit{"b", Interval{0,10}})
	vl.Add(Visit{"c", Interval{5,10}})
	
	
	
}

func equals(u, v VisitLogList) bool {
	for i := range u.visits {
		if u.visits[i] != v.visits[i] {
			return false
		}
	}
	return true
}
