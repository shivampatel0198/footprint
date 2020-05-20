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

func simulate(n *Node, i int, record *[][]NodeRecord) {
	for t := 0; t < iters; t++ {
		n.Walk()
		pos := n.Locate()
		n.Log(pos, t)
		if n.Infected {
			n.Push()
		} else {
			n.Check(func(overlaps map[PointCode][]Interval) bool {
				count := 0
				for _, intervals := range overlaps {
					for _, interval := range intervals {
						count += interval.Size()
					}
				} 
				return count > 1
			})

			// // Simple infection model: one-touch transmission
			// n.Check(func(overlaps map[PointCode][]Interval) bool {
			// 	return len(overlaps) > 0
			// })
		}
		r := NodeRecord{n.Id, pos, n.Infected}
		(*record)[t][i] = r
	}
}

type NodeRecord struct {
	NodeID   string
	Loc      Point
	Infected bool
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
	nodes[0] = NewNode("infected", NewSegmentedWalk(RandomPoint(0, node_spread)), g)
	nodes[0].MarkInfected()
	for i := 1; i < node_count; i++ {
		nodes[i] = NewNode("node "+strconv.Itoa(i), NewSegmentedWalk(RandomPoint(0, node_spread)), g)
	}

	// Simulate nodes
	for i, n := range nodes {
		simulate(n, i, &record)
	}

	// Filter out infected nodes
	infected := make([]string, 0)
	for _, n := range nodes {
		if n.Infected {
			infected = append(infected, n.Id)
		}
	}
	fmt.Println(infected)
	
	// Write JSON
	b, _ := json.Marshal(record)
	fmt.Println(string(b))
}
