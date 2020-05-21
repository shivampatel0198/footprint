package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

var (
	node_count  int
	node_spread int
	iters       int
)

/* Transmission Heuristic functions */

// Returns a threshold-checking function parametrized by k
func threshold(k int) func(map[PointCode][]Interval) bool {
	return func(overlaps map[PointCode][]Interval) bool {
		count := 0
		for _, intervals := range overlaps {
			for _, interval := range intervals {
				count += interval.Size()
			}
		}
		return count > k
	}
}

// Simple infection model: one-touch transmission
func one_touch() func(map[PointCode][]Interval) bool {
	return func(overlaps map[PointCode][]Interval) bool {
		return len(overlaps) > 0
	}
}

// Simulate one time step
func simulate(t Time, i int, n *Node, record *[][]NodeRecord) {

	// Move and log location
	p := n.Locate()
	n.Log(p, t)
	n.Walk()

	// Push/Check
	if n.Infected {
		n.Push()
	} else {
		n.Check(one_touch())
	}

	// Compute offset
	var offset Point
	if t == 0 {
		offset = Point{0, 0}
	} else {
		q := (*record)[t-1][i].Loc // previous pos
		offset = p.Sub(q)
	}

	(*record)[t][i] = NodeRecord{n.Id, p, offset, n.Infected}
}

type NodeRecord struct {
	NodeID   string
	Loc      Point
	Offset   Point
	Infected bool
}

// Returns a matrix displaying how many infected people have visited each position
func heatmap(record [][]NodeRecord) (out map[Point]int) {
	out = make(map[Point]int)
	for _, snapshot := range record {
		for _, node := range snapshot {
			if node.Infected {
				out[node.Loc]++
			}
		}
	}
	return
}

// Translate args into ints
func arg(i int) (out int) {
	out, _ = strconv.Atoi(os.Args[i])
	return
}

// Run the simulation
func main() {
	// Interpret command line args
	if len(os.Args) == 4 {
		node_count = arg(1)
		node_spread = arg(2)
		iters = arg(3)
	} else {
		fmt.Println("Usage: simulation [node-count] [node-spread] [iterations]")
		return
	}

	// Setup record
	record := make([][]NodeRecord, iters)
	for i := 0; i < iters; i++ {
		record[i] = make([]NodeRecord, node_count)
	}

	// Setup nodes
	nodes := make([]*Node, node_count)
	g := NewGlobalTrace()
	nodes[0] = NewNode(
		"infected",
		NewSegmentedWalk(RandomPoint(0, node_spread)),
		g,
	)
	nodes[0].MarkInfected()
	for i := 1; i < node_count; i++ {
		nodes[i] = NewNode(
			"node "+strconv.Itoa(i),
			NewSegmentedWalk(RandomPoint(0, node_spread)),
			g,
		)
	}

	// Simulate nodes
	for t := 0; t < iters; t++ {
		for i, n := range nodes {
			simulate(t, i, n, &record)
		}
	}

	// Filter out infected nodes
	// infected := make([]string, 0)
	// for _, n := range nodes {
	// 	if n.Infected {
	// 		infected = append(infected, n.Id)
	// 	}
	// }
	// fmt.Println(infected)

	// Write JSON
	b, _ := json.Marshal(record)
	fmt.Println(string(b))
}
