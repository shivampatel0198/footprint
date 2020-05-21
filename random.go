package main

import (
	"math/rand"
)

func Random(min, max int) int {
	// sign := rand.Intn(2)*2 - 1
	return min + rand.Intn(max-min+1)
}

func RandomExp() int {
	return int(rand.ExpFloat64())
}

func ChooseRandom(points []Point) Point {
	return points[rand.Intn(len(points))]
}

func RandomPoint(min, max int) Point {
	return Point{Random(min, max), Random(min, max)}
}
