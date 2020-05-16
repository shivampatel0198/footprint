package main

import (
	"fmt"
	"time"
	"strconv"
	"crypto/sha256"
)

var DHT = make(map[Point]VisitLogList)

type Node struct {
	key string                    // Private key used to generate tempIds
	log map[Point]VisitLogList    // Local visit log
}

func NewNode(key string) (n *Node) {
	n = new(Node)
	n.key = key
	n.log = make(map[Point]VisitLogList)
	return
}

func (n Node) String() string {
	return fmt.Sprintf("Key:%s\nLog:%v", n.key, n.log)
}

func (n Node) tempId() string {
	key := []byte(n.key + strconv.Itoa(time.Now().Nanosecond()))
	return fmt.Sprintf("%x", sha256.Sum256(key))
}

func (n Node) Log(cell Point, t Time) {
	for _, neighbor := range cell.ClosedNeighborhood() {
		vl := n.log[neighbor]
		if vl.Size() > 0 && t - vl.Last().Duration.End == TIME_STEP {
			n.log[neighbor] = extend(vl, t)
		} else {
			id := n.tempId()
			n.log[neighbor] = vl.Add(Visit{id, Interval{t, t}})
		}
	}
}

func (n Node) Push() {
	for p, vl := range n.log {
		ul, ok := DHT[p]
		if !ok {
			DHT[p] = vl
		} else {
			for _, visit := range vl.visits {
				DHT[p] = ul.Add(visit)
			}
		}
	}
}

func extend(vl VisitLogList, t Time) VisitLogList {
	return VisitLogList{
		append(
			vl.visits[:len(vl.visits)-1],
			Visit{vl.Last().Id, vl.Last().Duration.Extend(t)},
		),
	}
}
