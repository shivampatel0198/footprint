package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

const node_count = 50
const iters = 100

func simulate(n *Node, i int, record *[][]NodeRecord) {
	for t := 0; t < iters; t++ {
		pos := n.Locate()
		n.Log(pos, t)
		(*record)[t][i] = NodeRecord{n.Id, pos, n.Infected}
	}
	n.Check()
}

type NodeRecord struct {
	NodeID   string
	Loc      Point
	Infected bool
}

// Run the simulation
func main() {
	nodes := make([]*Node, node_count)
	g := NewGlobalTrace()
	path := []Point{
		Point{0, 0},
	}

	record := make([][]NodeRecord, iters)
	for i := 0; i < iters; i++ {
		record[i] = make([]NodeRecord, node_count)
	}

	// Setup infected node
	nodes[0] = NewNode("infected", NewCannedWalk(path), g)
	for t := 0; t < iters; t++ {
		pos := nodes[0].Locate()
		nodes[0].Log(pos, t)
		record[t][0] = NodeRecord{"infected", pos, true}
	}
	nodes[0].MarkInfected()

	// Setup and run other nodes
	for i := 1; i < node_count; i++ {
		nodes[i] = NewNode("node "+strconv.Itoa(i), NewRandomWalk(RandomPoint(0, 10)), g)
		simulate(nodes[i], i, &record)
	}

	b, _ := json.Marshal(record)
	fmt.Println(string(b))
}
