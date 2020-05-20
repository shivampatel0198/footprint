package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
)

var node_count = 50
var iters = 100

func simulate(n *Node, i int, record *[][]NodeRecord) {
	for t := 0; t < iters; t++ {
		pos := n.Locate()
		n.Log(pos, t)
		n.Check()
		(*record)[t][i] = NodeRecord{n.Id, pos, n.Infected}
	}
}

type NodeRecord struct {
	NodeID   string
	Loc      Point
	Infected bool
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

// Run the simulation
func main() {

	// Translate args into ints
	arg := func(i int) (out int) {
		out, _ = strconv.Atoi(os.Args[i])
		return
	}

	// Interpret command line args
	if len(os.Args) == 3 {
		node_count = arg(2)
		iters = arg(1)
	} else {
		fmt.Println("Usage: simulation [node-count] [iterations]")
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

	// Setup and run other nodes
	for i := 1; i < node_count; i++ {
		nodes[i] = NewNode("node "+strconv.Itoa(i), NewRandomWalk(RandomPoint(0, 10)), g)
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
