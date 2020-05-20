// Generates JSON objects storing randomly generated location data for nodes

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
)

type Point struct {
	X, Y int
}

func random(min, max int) int {
	sign := rand.Intn(2)*2 - 1
	return sign * (min + rand.Intn(max-min+1))
}

func randomStep(p Point, min, max int) Point {
	dx, dy := random(min, max), random(min, max)
	return Point{p.X + dx, p.Y + dy}
}

func WriteToFile(s string) {
	f, err := os.OpenFile("data.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(s); err != nil {
		log.Println(err)
	}
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: locations [min-step-size] [max-step-size] [node-count] [iterations]")
		return
	}

	arg := func(i int) (out int) {
		out, _ = strconv.Atoi(os.Args[i])
		return
	}

	min := arg(1)
	max := arg(2)
	nodeCount := arg(3)
	iters := arg(4)

	ps := make([][]Point, iters)

	// Randomly seed points
	for i := 0; i < nodeCount; i++ {
		ps[0] = make([]Point, nodeCount)
		ps[0][i] = Point{random(0, 5), random(0, 5)}
	}
	// Perform random walk
	for i := 1; i < iters; i++ {
		ps[i] = make([]Point, nodeCount)
		for j := range ps[i] {
			ps[i][j] = randomStep(ps[i-1][j], min, max)
		}
	}
	b, _ := json.Marshal(ps)
	WriteToFile(string(b))
}
