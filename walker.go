// Generates JSON objects storing randomly generated location data for nodes

package main

import (
	"math/rand"
)

type Walker interface {
	// Returns the current location in a random walk
	Where() Point

	// Moves to another location (or stays in place)
	Walk() 
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
	return Point{random(min,max), random(min,max)}
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
}

// Reads a canned walk from file.
func NewCannedWalk(ps []Point) (cw *CannedWalk) {
	cw = new(CannedWalk)
	cw.data = make([]Point, len(ps))
	for i := 0; i < len(ps); i++ {
		cw.data[i] = ps[i].Copy()
	}
	return
}

func (cw *CannedWalk) Walk() {
	cw.index++
}

func (cw *CannedWalk) Where() Point {
	l := len(cw.data)
	x := cw.index % (2*l - 2)
	if x/l == 0 {
		return cw.data[x]
	} else {
		return cw.data[2*(l-1)-x]
	}
}
