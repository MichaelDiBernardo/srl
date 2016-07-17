package math

import (
    "fmt"
)

// A point in 2D space.
type Point struct {
	X, Y int
}

// Stringifies a point.
func (p Point) String() string {
	return fmt.Sprintf("(%d,%d)", p.X, p.Y)
}

// Returns p + q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Returns p - q
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// (0,0). Also serves as the zero value for Point in comparisons.
var Origin Point
