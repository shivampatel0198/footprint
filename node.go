package main

import (
	"crypto/sha256"
	"fmt"
	"strconv"
	"time"
)

type Node struct {
	Id       string
	log      *LocalTrace
	bulletin *GlobalTrace
	gps      Walker
}

func NewNode(id string, bulletin *GlobalTrace) (n *Node) {
	n = new(Node)
	n.Id = id
	n.log = NewLocalTrace()
	n.bulletin = bulletin
	n.gps = new(RandomWalk)
	return
}

func (n Node) String() string {
	return fmt.Sprintf("Node %v\nLog:%v\nBulletin:%v", n.Id, n.log, &n.bulletin)
}

func TempId(key string) string {
	return fmt.Sprintf("%x",
		sha256.Sum256(
			[]byte(key+strconv.Itoa(time.Now().Nanosecond())),
		))
}

func Encode(cell Point) PointCode {
	return PointCode(fmt.Sprintf("%x",
		(sha256.Sum256(
			[]byte(cell.String()),
		))))
}

// Return the node's "current" cell location
func (n Node) Locate() Point {
	return n.gps.Where()
}

// Record local location information
func (n Node) Log(cell Point, t Time) {
	n.log.Add(Encode(cell), t)
}

// After testing positive, send local location information to the global store
func (n Node) Push() {
	n.log.Iterate(func(c PointCode, xs []Interval) {
		n.bulletin.Add(c, xs)
	})
}

/**
Determine contact events with infected individuals.

Looks at bulletin, looks at n's logged activity, and finds all of the intersections.
Returns a map from Points to Intervals containing all intersected intervals,
excluding those intervals belonging to n itself.
*/
func (n Node) Check() map[PointCode][]Interval {
	// (1) Pull needed cells from distributed hashtable
	// (2) Compute intersections
	return n.log.Intersect(n.bulletin)
}
