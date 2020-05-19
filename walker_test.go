package main

import (
	"testing"
)

func TestCannedWalk(t *testing.T) {
	path := []Point{
		Point{0, 0},
		Point{0, 1},
		Point{0, 2},
	}
	cw := NewCannedWalk(path)

	expected := []Point{
		Point{0, 0},
		Point{0, 1},
		Point{0, 2},
		Point{0, 2},
		Point{0, 1},
		Point{0, 0},
		Point{0, 0},
		Point{0, 1},
		Point{0, 2},
		Point{0, 2},
		Point{0, 1},
		Point{0, 0},
	}
	actual := make([]Point, len(expected))
	for i := range expected {
		actual[i] = cw.Where()
	}
	for i, e := range expected {
		if actual[i] != e {
			t.Errorf("expected=%v, actual %v", expected, actual)
		}
	}
}
