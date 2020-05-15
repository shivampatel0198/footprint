// Generates JSON objects storing randomly generated location data for nodes

package main

import (
	"fmt"
	"os"
	"strconv"
	"math/rand"
	"encoding/json"
)

type Point struct{
	X int
	Y int
}

func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

func Abs(x int) int {
	if x < 0 {
		return -x
	} else {
		return x
	}
}

// Manhattan distance
func Distance(u, v Point) int {
	return Abs(u.X - v.X) + Abs(u.Y - v.Y)
}

func RandomStep(x, min, max int ) int {
	sign := rand.Intn(2)*2 - 1
	return sign * (min + rand.Intn(max - min + 1))
}

func RandomWalk(p Point, min, max int) Point {
	dx, dy := RandomStep(p.X, min, max), RandomStep(p.Y, min, max)
	return Point{p.X + dx, p.Y + dy}
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: locations [min step size] [max step size] [node count] [iterations]")
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
	
	nodes := make([][]Point, iters)
	nodes[0] = make([]Point, nodeCount) //first row is all at origin
	for i := 1; i<iters; i++ {
		nodes[i] = make([]Point, nodeCount)
		for j := range nodes[i] {
			nodes[i][j] = RandomWalk(nodes[i-1][j], min, max)
		}
	}
	b, _ := json.Marshal(nodes)
	fmt.Printf("%s\n", b)

	for t := 1; t<len(nodes); t++ {
		for i := 0; i<len(nodes[0]); i++ {
			for j:= i+1; j<len(nodes[0]); j++ {
				if Distance(nodes[t][i], nodes[t][j]) < 2 {
					fmt.Printf("Nodes %d, %d are close at time %d: %s-%s\n", 
						i, j, t, nodes[t][i], nodes[t][j])
				}
			}
		}
	}
}
