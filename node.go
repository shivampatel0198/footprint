package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

// For debugging purposes.  In production build, this should be true.
const ENCODE = false

type Node struct {
	Id       string
	Infected bool
	log      *LocalTrace
	bulletin *GlobalTrace
	loc      Walker
}

func NewNode(id string, w Walker, bulletin *GlobalTrace) (n *Node) {
	n = new(Node)
	n.Id = id
	n.Infected = false
	n.log = NewLocalTrace()
	n.bulletin = bulletin
	n.loc = w
	return
}

func (n *Node) String() string {
	return fmt.Sprintf("Node %v\nLog:%v\nBulletin:%v", n.Id, n.log, &n.bulletin)
}

func TempId(key string) string {
	return fmt.Sprintf("%x",
		sha256.Sum256(
			[]byte(key+strconv.Itoa(time.Now().Nanosecond())),
		))
}

func Encode(cell Point) PointCode {
	if ENCODE {
		return PointCode(fmt.Sprintf("%x",
			(sha256.Sum256(
				[]byte(cell.String()),
			))))
	} else {
		return PointCode(cell.String())
	}
}

func (n *Node) Walk() {
	n.loc.Walk()
}

// Return the node's "current" cell location
func (n *Node) Locate() Point {
	return n.loc.Where()
}

// Record local location information
func (n *Node) Log(cell Point, t Time) {
	n.log.Add(Encode(cell), t)
}

func (n *Node) MarkInfected() {
	n.Infected = true
}

// After testing positive, send local location information to the global store
func (n *Node) Push() {
	n.log.Iterate(func(c PointCode, xs []Interval) {
		n.bulletin.Add(c, xs)
	})
}

/**
Determine contact events with infected individuals.

Looks at bulletin, looks at n's logged activity, and finds all of the intersections.
Using those intersections, determine whether n is now infected/at risk.  
*/
func (n *Node) Check(pred func(overlaps map[PointCode][]Interval) bool) {
	// Ignore already infected nodes
	if n.Infected {
		return
	}
	if pred(n.log.Intersect(n.bulletin)) {
		n.MarkInfected()
	}
}
