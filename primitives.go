package main

import (
	"fmt"
)

const (
	TIME_STEP    int = 1
	EXPOS_THRESH int = 3
	DIST_THRESH  int = 2
)

type Point struct {
	X, Y int
}

func (p Point) Copy() Point {
	return Point{p.X, p.Y}
}

func (p Point) Add(u Point) Point {
	return Point{p.X + u.X, p.Y + u.Y}
}

func (p Point) Sub(u Point) Point {
	return Point{p.X - u.X, p.Y - u.Y}
}

func (p Point) ClosedNeighborhood() (out []Point) {
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			out = append(out, Point{p.X + dx, p.Y + dy})
		}
	}
	return
}

func (p Point) String() string {
	return fmt.Sprintf("(%v,%v)", p.X, p.Y)
}

type Time = int

func Max(a, b Time) Time {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a, b Time) Time {
	if a < b {
		return a
	} else {
		return b
	}
}
