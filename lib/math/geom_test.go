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
	p1, p2, want := Pt(1, 2), Pt(3, 5), Pt(4, 7)
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
	p1, p2, want := Pt(1, 2), Pt(3, 5), Pt(-2, -3)
	diff := p1.Sub(p2)

	if diff != want {
		t.Errorf("%v.Sub(%v) was %v, want %v", p1, p2, diff, want)
	}

	want = Pt(2, 3)
	diff = p2.Sub(p1)

	if diff != want {
		t.Errorf("%v.Sub(%v) was %v, want %v", p2, p1, diff, want)
	}
}

func TestRectDimensions(t *testing.T) {
	sut := Rect(Pt(2, 3), Pt(11, 5))
	if expected, actual := 9, sut.Width(); expected != actual {
		t.Errorf("Rectangle.Width() is %d; want %d", actual, expected)
	}
	if expected, actual := 2, sut.Height(); expected != actual {
		t.Errorf("Rectangle.Height() is %d; want %d", actual, expected)
	}
}

func TestRectHasPoint(t *testing.T) {
	sut := Rect(Pt(1, 2), Pt(4, 5))
	if pt := Pt(1, 3); !sut.HasPoint(pt) {
		t.Errorf("%v.HasPoint(%v) is false; want true", sut, pt)
	}

	if pt := Pt(4, 5); sut.HasPoint(pt) {
		t.Errorf("%v.HasPoint(%v) is true; want false", sut, pt)
	}
}
