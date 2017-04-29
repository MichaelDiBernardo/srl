package math

import (
	"fmt"
	"math"
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

// "Filters" the given list of points into a new list that only contains those
// points which are contained in r.
func (r Rectangle) Clip(pts []Point) []Point {
	clipped := make([]Point, 0, len(pts))
	for _, pt := range pts {
		if pt.In(r) {
			clipped = append(clipped, pt)
		}
	}
	return clipped
}

// Yields rectangle that contains all coordinates of Chebyshev distance <= r
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

// Creates a list of points that have Chebyshev distance == r from (0, 0).  The
// points are ordered, following them will take you from topleft -> bottomleft
// -> bottomright -> topright.
func ChebyEdge(r int) []Point {
	if r <= 0 {
		return []Point{Pt(0, 0)}
	}

	edge := make([]Point, 0)

	for i := -r; i < r; i++ {
		edge = append(edge, Pt(i, -r))
	}
	for i := -r; i < r; i++ {
		edge = append(edge, Pt(r, i))
	}
	for i := -r; i < r; i++ {
		edge = append(edge, Pt(-i, r))
	}
	for i := -r; i < r; i++ {
		edge = append(edge, Pt(-r, -i))
	}

	return edge
}

// Calculates the Chebyshev distance between two points.
func ChebyDist(p1, p2 Point) int {
	return Max(Abs(p1.X-p2.X), Abs(p1.Y-p2.Y))
}

// The euclidean distance between two points, floored to the nearest integer.
func EucDist(p1, p2 Point) int {
	diffxs := math.Pow(float64(p1.X-p2.X), 2)
	diffys := math.Pow(float64(p1.Y-p2.Y), 2)
	return int(math.Sqrt(diffxs + diffys))
}

// Get all points adjacent to pt.
func Adj(pt Point) []Point {
	around := ChebyEdge(1)
	adj := make([]Point, 0, 8)
	for _, p := range around {
		adj = append(adj, pt.Add(p))
	}
	return adj
}

// Approximate a straight line between p1 and p2.
// Uses Bresenham's algorithm, stolen and adapted from
// http://www.roguebasin.com/index.php?title=Bresenham%27s_Line_Algorithm#Go
func Line(p1, p2 Point) (line []Point) {
	x1, y1 := p1.X, p1.Y
	x2, y2 := p2.X, p2.Y

	isSteep := Abs(y2-y1) > Abs(x2-x1)
	if isSteep {
		x1, y1 = y1, x1
		x2, y2 = y2, x2
	}

	reversed := false
	if x1 > x2 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
		reversed = true
	}

	deltaX := x2 - x1
	deltaY := Abs(y2 - y1)
	err := deltaX / 2
	y := y1
	var ystep int

	if y1 < y2 {
		ystep = 1
	} else {
		ystep = -1
	}

	for x := x1; x < x2+1; x++ {
		if isSteep {
			line = append(line, Pt(y, x))
		} else {
			line = append(line, Pt(x, y))
		}
		err -= deltaY
		if err < 0 {
			y += ystep
			err += deltaX
		}
	}

	if reversed {
		for i, j := 0, len(line)-1; i < j; i, j = i+1, j-1 {
			line[i], line[j] = line[j], line[i]
		}
	}

	return
}
