package main

func abs(val int) int {
	if (val < 0) {
		return val * -1
	}
	return val
}

func copysign(mag int, sign int) int {
	if ((sign < 0) && (mag > 0)) ||
	   ((sign > 0) && (mag < 0)) {
		return mag * -1
	}
	return mag
}