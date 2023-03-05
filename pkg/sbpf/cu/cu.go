package cu

// This file contains helper routines for the calculation of compute units.

func ConsumeLowerBound(cu int, lower int, x int) int {
	if x < lower {
		return cu - lower
	}
	return cu - x
}
