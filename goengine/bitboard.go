package goengine

const debruijn64 uint64 = 0x03f79d71b4cb0a89
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

var rayAttacks [64][8]uint64

type RayDir uint8
const (
	NORTH RayDir = iota
	N_EAST
	EAST
	S_EAST
	SOUTH
	S_WEST
	WEST
	N_WEST
)

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

func initRayAttacks() {
	// Calculate ray attacks
	// TODO: Find a more elegant approach to ray-move calculation
	for i, _ := range rayAttacks {
		row := i / 8
		col := i % 8  
		
		// Calculate north ray attacks
		for j := 8; j > row; j-- {
			rayAttacks[i][NORTH] = moveNorth((1 << i) | rayAttacks[i][NORTH])
		}

		// Calculate north-east ray attacks
		for j := 8; j > row; j-- {
			rayAttacks[i][N_EAST] = moveNEast((1 << i) | rayAttacks[i][N_EAST])
		}

		// Calculate east ray attacks
		for j := 0; j < col; j++ {
			rayAttacks[i][EAST] = moveEast((1 << i) | rayAttacks[i][EAST])
		}

		// Calculate south-east ray attacks
		for j := row; j > 0; j-- {
			rayAttacks[i][S_EAST] = moveSEast((1 << i) | rayAttacks[i][S_EAST])
		}

		// Calculate south ray attacks
		for j := row; j > 0; j-- {
			rayAttacks[i][SOUTH] = moveSouth((1 << i) | rayAttacks[i][SOUTH])
		}

		// Calculate south-west ray attacks
		for j := row; j > 0; j-- {
			rayAttacks[i][S_WEST] = moveSWest((1 << i) | rayAttacks[i][S_WEST])
		}

		// Calculate west ray attacks
		for j := col; j < 8; j++ {
			rayAttacks[i][WEST] = moveWest((1 << i) | rayAttacks[i][WEST])
		}

		// Calculate north-west ray attacks
		for j := 8; j > row; j-- {
			rayAttacks[i][N_WEST] = moveNWest((1 << i) | rayAttacks[i][N_WEST])
		}
	}
}

func getTransSet(piece uint64, occup uint64) uint64 {
	var set uint64 = 0
	for piece != 0 {
		var sqr uint8 = bitScanForward(piece)
		set |= getPosRayAttacks(sqr, occup, NORTH)
		set |= getNegRayAttacks(sqr, occup, EAST)
		set |= getPosRayAttacks(sqr, occup, WEST)
		set |= getNegRayAttacks(sqr, occup, SOUTH)
		piece ^= (1 << sqr)
	}
	return set
}

func getDiagSet(bb uint64, occup uint64) uint64 {
	var set uint64 = 0
	for bb != 0 {
		var sqr uint8 = bitScanForward(bb)
		set |= getPosRayAttacks(sqr, occup, N_EAST)
		set |= getPosRayAttacks(sqr, occup, N_WEST)
		set |= getNegRayAttacks(sqr, occup, S_EAST)
		set |= getNegRayAttacks(sqr, occup, S_WEST)
		bb ^= (1 << sqr)
	}
	return set
}

func getPosRayAttacks(sqr uint8, occup uint64, dir RayDir) uint64 {
	var attacks uint64 = rayAttacks[sqr][dir]
	var blockers uint64 = attacks & (^occup)
	sqr = bitScanForward(blockers | (0x8000000000000000))
	return attacks ^ rayAttacks[sqr][dir]
}

func getNegRayAttacks(sqr uint8, occup uint64, dir RayDir) uint64 {
	var attacks uint64 = rayAttacks[sqr][dir]
	var blockers uint64 = attacks & (^occup)
	sqr = bitScanReverse(blockers | 1)
	return attacks ^ rayAttacks[sqr][dir]
}

func moveNWest(board uint64) uint64 {return (board << 9) & (NOT_H_FILE)}
func moveNorth(board uint64) uint64 {return board << 8}
func moveNEast(board uint64) uint64 {return (board << 7) & (NOT_A_FILE)}
func moveEast(board uint64) uint64 {return (board >> 1) & (NOT_A_FILE)}
func moveSEast(board uint64) uint64 {return (board >> 9) & (NOT_A_FILE)}
func moveSouth(board uint64) uint64 {return board >> 8}
func moveSWest(board uint64) uint64 {return (board >> 7) & (NOT_H_FILE)}
func moveWest(board uint64) uint64 {return (board << 1) & (NOT_H_FILE)}