// Generates JSON objects storing randomly generated location data for nodes

package main

import (
	"math/rand"
)

type Walker interface {
	// Returns the next location in a random walk
	Where() Point
}

// Walk randomly on 2D grid
type RandomWalk Point

func NewRandomWalk(start Point) *RandomWalk {
	rw := new(RandomWalk)
	rw.set(start)
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

func (rw *RandomWalk) set(p Point) {
	*rw = RandomWalk(p)
}

func (rw *RandomWalk) Where() Point {
	defer rw.set(randomStep(Point(*rw)))
	return Point(*rw)
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

func (cw *CannedWalk) Inc() {
	cw.index++
}

func (cw *CannedWalk) Where() Point {
	defer cw.Inc()

	i, l, x := 0, len(cw.data), cw.index
	if x/l%2 == 0 {
		i = x % l
	} else {
		i = l - (x % l) - 1
	}
	return cw.data[i]
}
