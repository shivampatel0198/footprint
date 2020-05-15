package main

import (
	"fmt"
	"time"
	"strconv"
	"crypto/sha256"
)

var DHT = make(map[Point]VisitLog)

type Node struct {
	key string                    // Private key used to generate tempIds
	ids map[string]bool           // Set of previously exposed tempIds
	log map[Point][]Interval      // Local visit log
}

func NewNode(key string) (n *Node) {
	n = new(Node)
	n.key = key
	n.ids = make(map[string]bool)
	n.log = make(map[Point][]Interval)
	return
}

func (n Node) String() string {
	return fmt.Sprintf("Key:%s\nExposed-IDs:%v\nLog:%v", n.key, n.ids, n.log)
}

func (n Node) tempId() string {
	key := []byte(n.key + strconv.Itoa(time.Now().Nanosecond()))
	return fmt.Sprintf("%x", sha256.Sum256(key))
}

// TODO: make more efficient
func (n Node) Log(cell Point, t Time) {
	neighbors := cell.ClosedNeighborhood()
	for _, neighbor := range neighbors {
		l := n.log[neighbor]
		if len(l) > 0 && t - l[len(l)-1].End == TIME_STEP {
			n.log[neighbor][len(l)-1] = l[len(l)-1].Extend(t)
		} else {
			n.log[neighbor] = append(l, Interval{t, t})
		}
	}
}

func (n Node) Push() {
	for cell, intervals := range n.log {

		// Initialize missing cells
		_, ok := DHT[cell]
		if !ok {
			DHT[cell] = VisitLogList{}
		}

		// Push intervals
		for _, interval := range intervals {
			id := n.tempId()
			n.ids[id] = true //track emitted tempIds
			DHT[cell] = DHT[cell].Add(Visit{id, interval})
		}
	}
}

func main() {
	var n = NewNode("katya")
	fmt.Println(n.tempId())

	n.Log(Point{0,0}, 0)
	n.Log(Point{0,1}, 1)
	n.Log(Point{1,0}, 2)
	n.Log(Point{1,1}, 3)
	n.Log(Point{1,1}, 4)
	n.Log(Point{1,0}, 5)
	n.Log(Point{0,1}, 6)
	n.Log(Point{0,0}, 7)

	n.Push()

	fmt.Println(*n)
	fmt.Println("DHT\n", DHT)
}
