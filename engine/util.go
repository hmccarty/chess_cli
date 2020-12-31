package engine

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

var index64 = [64]uint8 {
	0, 47,  1, 56, 48, 27,  2, 60,
   57, 49, 41, 37, 28, 16,  3, 61,
   54, 58, 35, 52, 50, 42, 21, 44,
   38, 32, 29, 23, 17, 11,  4, 62,
   46, 55, 26, 59, 40, 36, 15, 53,
   34, 51, 20, 43, 31, 22, 10, 45,
   25, 39, 14, 33, 19, 30,  9, 24,
   13, 18,  8, 12,  7,  6,  5, 63,
}
const debruijn64 uint64 = 0x03f79d71b4cb0a89

func bitScanForward(board uint64) uint8 {
	return index64[((board ^ (board-1)) * debruijn64) >> 58]
}

func bitScanReverse(board uint64) uint8 {
	board |= board >> 1
	board |= board >> 2
	board |= board >> 4
	board |= board >> 8
	board |= board >> 16
	board |= board >> 32
	return index64[(board * debruijn64) >> 58]
}