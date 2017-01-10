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

// A rectangle with top-left point Min and top-right point Max, such that it
// contains points where with Min.X <= p.X < Max.X and Min.Y <= p.Y < Max.Y.
type Rectangle struct {
	Min, Max Point
}

var ZeroRect Rectangle

func Rect(min, max Point) Rectangle {
	return Rectangle{Min: min, Max: max}
}

func (r Rectangle) Width() int {
	return r.Max.X - r.Min.X
}

func (r Rectangle) Height() int {
	return r.Max.Y - r.Min.Y
}

func (r Rectangle) HasPoint(p Point) bool {
	return r.Min.X <= p.X && r.Min.Y <= p.Y && p.X < r.Max.X && p.Y < r.Max.Y
}

func (r Rectangle) String() string {
	return fmt.Sprintf("Rect(%v, %v)", r.Min, r.Max)
}

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
// Stolen from https://golang.org/src/image/geom.go?m=text.
func (r Rectangle) Intersect(s Rectangle) Rectangle {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	if r.Min.X > r.Max.X || r.Min.Y > r.Max.Y {
		return ZeroRect
	}
	return r
}

func (r Rectangle) Center() Point {
	return Pt((r.Min.X+r.Max.X)/2, (r.Min.Y+r.Max.Y)/2)
}

// Yields rectangle that contains all coordinates of Chebyshev distance 'r'
// from (0, 0). You can iterate over all points like:
//
// points := Chebyshev(3)
// for y := points.Min.Y; y < points.Max.Y; y++ {
//   for x := points.Min.X; x < points.Max.X; x++ {
//       -- do something with Pt(x, y)
//   }
// }
func Chebyshev(r int) Rectangle {
	return Rect(Pt(-r, -r), Pt(r+1, r+1))
}
