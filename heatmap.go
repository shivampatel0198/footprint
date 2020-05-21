package main

import (
	"fmt"
	"math"
)

type HeatmapCell struct {
	X, Y  int
	Weight float64
}

// Given a snapshot, computes the number of infected nodes for each point
func count(record []NodeRecord) (out map[Point]int) {
	out = make(map[Point]int)
	for _, nr := range record {
		if nr.Infected {
			out[nr.Loc]++
		}
	}
	return
}

func counts(t int, record [][]NodeRecord) []map[Point]int {
	if len(record) <= t {
		panic(fmt.Sprintf("%v <= %v", len(record), t))
	}
	cs := make([]map[Point]int, t+1)
	for i, snapshot := range record[:t+1] {
		cs[i] = count(snapshot)
	}
	return cs
}

func weight(t, T int) float64 {
	return math.Pow(0.5, float64(T-t+1))
}

func weighted_sum(ms []map[Point]int) map[Point]float64 {
	out := make(map[Point]float64)
	for t, m := range ms {
		for p, x := range m {
			out[p] += float64(x) * weight(t, len(m))
		}
	}
	return out
}

func floatify(m map[Point]int) map[Point]float64 {
	out := make(map[Point]float64)
	for k, v := range m {
		out[k] = float64(v)
	}
	return out
}

// Reduces a slice of maps into a single map by summing all
// key-value pairs
func sum(ms []map[Point]int) map[Point]int {
	out := make(map[Point]int)
	for _, m := range ms {
		for p, x := range m {
			out[p] += x
		}
	}
	return out
}

// Maps all values in the input map to the interval [0,1]
// Note: Assumes that entries are non-negative
func normalize(m map[Point]float64) (n map[Point]float64) {

	// Copy m into n
	n = make(map[Point]float64)
	tot := 0.0
	for p, x := range m {
		n[p] = x
		tot += x
	}

	// Avoid div by 0
	if tot == 0 {
		tot = 1
	}

	// Normalize
	for p := range n {
		n[p] /= tot
	}
	return
}

// Returns a single heatmap for time t
func heatmap(t int, record [][]NodeRecord) map[Point]float64 {
	return normalize(floatify(sum(counts(t, record))))
}

// Returns all heatmaps for a record
func Heatmaps(record [][]NodeRecord) [][]HeatmapCell {
	out := make([][]HeatmapCell, len(record))
	for t := range record {
		out[t] = flatten(heatmap(t, record))
	}
	return out
}

// Heatmap where more recent visits are weighted more heavily
func weighted_heatmap(t int, record [][]NodeRecord) map[Point]float64 {
	return normalize(weighted_sum(counts(t, record)))
}

// Returns all cumulative heatmaps generated by record
func WeightedHeatmaps(record [][]NodeRecord) [][]HeatmapCell {
	out := make([][]HeatmapCell, len(record))
	for t := range record {
		out[t] = flatten(weighted_heatmap(t, record))
	}
	return out
}

func flatten(m map[Point]float64) []HeatmapCell {
	out := make([]HeatmapCell, 0)
	for k, v := range m {
		out = append(out, HeatmapCell{k.X, k.Y, v})
	}
	return out
}
