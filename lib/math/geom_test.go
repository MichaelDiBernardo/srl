package math

import (
	"testing"
)

func TestOriginWorksAsZeroValue(t *testing.T) {
	var zero Point

	if zero != Origin {
		t.Error("zero Point != math.Origin")
	}
}

func TestAddPoints(t *testing.T) {
	p1, p2, want := Point{1, 2}, Point{3, 5}, Point{4, 7}
	sum := p1.Add(p2)

	if sum != want {
		t.Errorf("%v.Add(%v) was %v, want %v", p1, p2, sum, want)
	}

	sum = p2.Add(p1)

	if sum != want {
		t.Errorf("%v.Add(%v) was %v, want %v", p2, p1, sum, want)
	}
}

func TestSubPoints(t *testing.T) {
	p1, p2, want := Point{1, 2}, Point{3, 5}, Point{-2, -3}
	diff := p1.Sub(p2)

	if diff != want {
		t.Errorf("%v.Sub(%v) was %v, want %v", p1, p2, diff, want)
	}

    want = Point{2, 3}
	diff = p2.Sub(p1)

	if diff != want {
		t.Errorf("%v.Sub(%v) was %v, want %v", p2, p1, diff, want)
	}
}
