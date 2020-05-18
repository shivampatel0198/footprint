package main

import (
	"fmt"
	"time"
	"strconv"
	"crypto/sha256"
)

type Node struct {
	trace *LocalTrace
	bulletin *GlobalTrace
}

func NewNode(bulletin *GlobalTrace) (n *Node) {
	n = new(Node)
	n.trace = NewLocalTrace()
	n.bulletin = bulletin
	return
}

func (n Node) String() string {
	return fmt.Sprintf("Trace:\n%v\nBulletin:%v", n.trace, &n.bulletin)
}

func TempId(key string) string {
	return fmt.Sprintf("%x",
		sha256.Sum256(
			[]byte(key + strconv.Itoa(time.Now().Nanosecond())),
		))
}

func Encode(cell Point) PointCode {
	return PointCode(fmt.Sprintf("%x",
		(sha256.Sum256(
			[]byte(cell.String()),
		))))
}

// Record local location information
func (n Node) Log(cell Point, t Time) {
	n.trace.Add(Encode(cell), t)
}

// After testing positive, send local location information to the global store
func (n Node) Push() {
	n.trace.Iterate(func(c PointCode, xs []Interval) {
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
	return n.trace.Intersect(n.bulletin)
}

