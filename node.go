package main

import (
	"fmt"
	"time"
	"strconv"
	"crypto/sha256"
)

type Node struct {
	key string                    // Private key used to generate tempIds
	ids map[string]bool           // Set of previously exposed tempIds
	log map[Point][]Interval      // Local visit log
	DHT *map[Point]VisitLog       // Reference to DHT
}

func NewNode(key string) (n *Node) {
	n = new(Node)
	n.key = key
	n.ids = make(map[string]bool)
	n.log = make(map[Point][]Interval)
	n.DHT = new(map[Point]VisitLog)
	return
}

func (n Node) String() string {
	return fmt.Sprintf("Key:%s\nExposed-IDs:%v\nLog:%v", n.key, n.ids, n.log)
}

func (n Node) tempId() string {
	key := []byte(n.key + strconv.Itoa(time.Now().Nanosecond()))
	return fmt.Sprintf("%x", sha256.Sum256(key))
}

/*
  if cell == node.prevCell:
    node.hasNotMoved = true
    return

  // Track the cells around cell that are close enough to the node
  // to constitute transmission events
  let cells = get-neighbors(cell) 

  for cell in cells:
    add interval(t) to node.log[cell], merging into existing intervals if possible,
    taking node.hasNotMoved into account
*/
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

func main() {
	var n = NewNode("katya")
	fmt.Println(n.tempId())
	fmt.Println(n.tempId())
	fmt.Println(n.tempId())
	n.Log(Point{0,0}, 0)
	n.Log(Point{0,1}, 1)
	n.Log(Point{1,0}, 2)
	n.Log(Point{1,1}, 3)
	fmt.Println(*n)
}
