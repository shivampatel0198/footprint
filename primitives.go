package main

import "fmt"

//Constants
var TIME_STEP = 1
var EXPOS_THRESH = 3
var DIST_THRESH = 2

type Point struct {
	X int
	Y int
}

func (p Point) ClosedNeighborhood() (out []Point) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			out = append(out, Point{p.X+dx, p.Y+dy})
		}
	}
	return
}

func (p Point) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

// Allow for substitution of different datatypes (float64, time.Time)
type Time = int

func max(a, b Time) Time {
	if a > b { 
		return a
	} else { 
		return b	
	}
}

func min(a, b Time) Time {
	if a < b {
		return a
	} else {
		return b
	}
}
