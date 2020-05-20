package main

import (
	"encoding/json"
	"fmt"
	"log"
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
		pos := n.Locate()
		n.Log(pos, t)
		if n.Infected {
			n.Push()
		} else {
			// Simple infection model: one-touch transmission
			n.Check(func(overlaps map[PointCode][]Interval) bool {
				return len(overlaps) > 0
			})
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

func WriteToFile(s string) {
	f, err := os.OpenFile("data.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(s); err != nil {
		log.Println(err)
	}
}

// Run the simulation
func main() {

	// Translate args into ints
	arg := func(i int) (out int) {
		out, _ = strconv.Atoi(os.Args[i])
		return
	}

	// Interpret command line args
	if len(os.Args) == 4 {
		node_count = arg(1)
		node_spread = arg(2)
		iters = arg(3)
	} else {
		fmt.Println("Usage: simulation [node-count] [node-spread] [iterations]")
		return
	}

	// Setup nodes
	nodes := make([]*Node, node_count)
	g := NewGlobalTrace()
	path := []Point{
		Point{0, 0},
	}

	// Setup record
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
	nodes[0].Push()

	// Setup and run other nodes
	for i := 1; i < node_count; i++ {
		nodes[i] = NewNode("node "+strconv.Itoa(i), NewRandomWalk(RandomPoint(0, node_spread)), g)
		simulate(nodes[i], i, &record)
	}

	infected := make([]string, 0)
	for _, n := range nodes {
		if n.Infected {
			infected = append(infected, n.Id)
		}
	}
	fmt.Println("Infected nodes:", infected)

	b, _ := json.Marshal(record)
	WriteToFile(string(b))
}
