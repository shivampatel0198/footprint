// Generates JSON objects storing randomly generated location data for nodes

package main

import (
	"math/rand"
)

type Walker interface {
	// Moves to another location (or stays in place)
	Walk()

	// Returns the current location in a random walk
	Where() Point
}

type StationaryWalk struct {
	pos Point
}

func NewStationaryWalk(start Point) *StationaryWalk {
	sw := new(StationaryWalk)
	sw.pos = start
	return sw
}

func (sw *StationaryWalk) Walk() {}

func (sw *StationaryWalk) Where() Point {
	return sw.pos
}

// Walk randomly on 2D grid
type RandomWalk struct {
	pos Point
}

func NewRandomWalk(start Point) *RandomWalk {
	rw := new(RandomWalk)
	rw.pos = start
	return rw
}

func random(min, max int) int {
	sign := rand.Intn(2)*2 - 1
	return sign * (min + rand.Intn(max-min+1))
}

func RandomPoint(min, max int) Point {
	return Point{random(min, max), random(min, max)}
}

// Take a random step: x,y in [-1, 1]
func randomStep(p Point) Point {
	dx, dy := random(0, 1), random(0, 1)
	return Point{p.X + dx, p.Y + dy}
}

func (rw *RandomWalk) Walk() {
	rw.pos = randomStep(rw.pos)
}

func (rw *RandomWalk) Where() Point {
	return rw.pos
}

// Read in data to walk.
// If data runs out, run backwards through it to start, then back forwards...
type CannedWalk struct {
	data  []Point
	index int
	dx    int
}

// Reads a canned walk from file.
func NewCannedWalk(ps []Point) (cw *CannedWalk) {
	cw = new(CannedWalk)
	cw.data = make([]Point, len(ps))
	for i := 0; i < len(ps); i++ {
		cw.data[i] = ps[i].Copy()
	}
	cw.dx = 1
	return
}

func (cw *CannedWalk) Walk() {
	if cw.index >= len(cw.data) {
		cw.index = len(cw.data)
		cw.dx = -1
	} else if cw.index < 0 {
		cw.index = -1
		cw.dx = 1
	}
	cw.index += cw.dx
}

func (cw *CannedWalk) Where() Point {
	return cw.data[cw.index]
}

// Repeatedly rolls two dice:
// one to choose walking direction, one to choose walking distance.
type SegmentedWalk struct {

	// Current position
	pos Point

	// Store direction as a displacement vector to be added at each time step
	direction Point

	// Stores remaining distance to walk before re-rolling
	distance int
}

func NewSegmentedWalk(start Point) *SegmentedWalk {
	sw := new(SegmentedWalk)
	sw.pos = start
	sw.reroll()
	return sw
}

func (sw *SegmentedWalk) reroll() {
	sw.distance = int(rand.ExpFloat64()) + 1
	sw.direction = Point{random(0, 1), random(0, 1)}
}

func (sw *SegmentedWalk) Walk() {
	if sw.distance == 0 {
		sw.reroll()
	}
	sw.pos = sw.pos.Add(sw.direction)
	sw.distance--
}

func (sw *SegmentedWalk) Where() Point {
	return sw.pos
}
