package main

import (
	"fmt"
	// "time"
	"strconv"
)

const node_count = 5
const run_time = 3

func simulate(n *Node, quit chan struct{}) {

	// Move around
	for t := 0; t < run_time; t++ {
		n.Log(n.Locate(), t)
	}

	// Periodically check against bulletin
	// go func() {
	// 	for {
	// 		select {
	// 		case <-quit:
	// 			return
	// 		default:
	// 			time.Sleep(100 * time.Millisecond)
	// 			overlaps := n.Check()
	// 			if len(overlaps) > 0 {
	// 				fmt.Println("Overlap!")
	// 				fmt.Println(n, overlaps)
	// 			}
	// 		}
	// 	}
	// }()
}

// Run the simulation
func main() {
	ns := make([]*Node, node_count)
	g := NewGlobalTrace()
	for i := 0; i < node_count; i++ {
		ns[i] = NewNode(strconv.Itoa(i), g)
	}
	for t := 0; t < run_time; t++ {
		ns[0].Log(ns[0].Locate(), t)
	}
	ns[0].Push()

	quit := make(chan struct{})
	for _, n := range ns {
		simulate(n, quit)
	}
	for _, n := range ns {
		fmt.Println(n)
	}
	// close(quit)
}
