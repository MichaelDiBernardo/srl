package math

func Min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func Max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func Sgn(x int) int {
	switch {
	case x < 0:
		return -1
	case x == 0:
		return 0
	default:
		return 1
	}
}

func Abs(x int) int {
	if x >= 0 {
		return x
	} else {
		return -1 * x
	}
}
