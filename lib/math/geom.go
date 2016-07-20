package math

import (
	"fmt"
)

// A point in 2D space.
type Point struct {
	X, Y int
}

// Shorthand to create a new point.
func Pt(x, y int) Point {
	return Point{x, y}
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

// Delegates to s.HasPoint -- exists because pt.In(s) is prettier to read than
// s.HasPoint(p) sometimes.
func (p Point) In(s Shape) bool {
	return s.HasPoint(p)
}

// (0,0). Also serves as the zero value for Point in comparisons.
var Origin Point

type Shape interface {
	HasPoint(p Point) bool
}
