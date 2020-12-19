package rtree

import (
	"math/rand"
	"sort"
)

// M is MaxSize
const M = 10e5

// Point represents a  xy
type Point struct {
	X  int
	Y  int
	ID int
}

// Query represents a 2d rectangle
type Query struct {
	A  *Point
	B  *Point
	ID int
}

func new(x, y int) *Point {
	return &Point{
		X: x,
		Y: y,
	}
}

func getRand(min, max int) int {
	return rand.Intn(max) + min
}

// GPoint returns the address of a new xy point
func GPoint(min, max int) *Point {
	return new(getRand(min, max), getRand(min, max))
}

// PointSet returns a set of x,y points
func PointSet(n int) []*Point {
	pointSet := make([]*Point, n, n)
	for i := 0; i < n; i++ {
		pointSet[i] = GPoint(1, M)
	}
	return pointSet
}

// GQuery returns a new query
func GQuery(s int) *Query {
	pnt := GPoint(1, M-s)
	return &Query{
		A: pnt,
		B: new(pnt.X+s, pnt.Y+s),
	}
}

// SortedPointSet returns a set of sortex random x,y coords
func SortedPointSet(n int) []*Point {
	pointSet := PointSet(n)
	sort.Slice(pointSet, func(i, j int) bool { return pointSet[i].X < pointSet[j].X })
	return pointSet
}
