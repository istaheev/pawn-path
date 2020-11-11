package main

import (
	"fmt"
	"sort"
)

type position struct {
	x int
	y int
}

// move point by specified offset
func (p position) move(dir position) position {
	return position{x: p.x + dir.x, y: p.y + dir.y}
}

type solver struct {
	// Board size and available directions (read-only)
	width      int
	height     int
	directions []position
	// Tracks visited cells
	visited []bool
	// A list of moves forming current path
	moves []position
	// Helper buffers to avoid excessive memory allocations
	movesBuf   []position
	movesCount []int
}

func newSolver(width, height int, directions []position) solver {
	var boardSize = width * height
	return solver{
		width:      width,
		height:     height,
		directions: directions,
		visited:    make([]bool, boardSize),
		moves:      make([]position, 0, boardSize),
		movesBuf:   make([]position, 0, len(directions)),
		movesCount: make([]int, len(directions)),
	}
}

// cellIndex returns a zero-based index of the cell to be used as an index in
// other arrays. Index is always in [0, boardSize) range.
func (s *solver) cellIndex(p position) int {
	return p.y*s.width + p.x
}

func (s *solver) inside(p position) bool {
	return p.x >= 0 && p.x < s.width && p.y >= 0 && p.y < s.height
}

func (s *solver) boardSize() int {
	return s.width * s.height
}

func (s *solver) isVisited(p position) bool {
	return s.visited[s.cellIndex(p)]
}

func (s *solver) setVisited(p position, visited bool) {
	s.visited[s.cellIndex(p)] = visited
}

func (s *solver) canJump(p position) bool {
	return s.inside(p) && !s.isVisited(p)
}

func (s *solver) jump(p position) bool {
	if !s.canJump(p) {
		return false
	}
	s.moves = append(s.moves, p)
	s.setVisited(p, true)
	return true
}

// undoMoves marks all the moves in the range [move..] as if they didn't happen
func (s *solver) undoMoves(move int) {
	for i := move; i < len(s.moves); i++ {
		var p = s.moves[i]
		s.setVisited(p, false)
	}
	s.moves = s.moves[0:move]
}

func (s *solver) findPath(start position) []position {
	// To avoid recursion the function uses its own heap-based stack.
	// An entry holds a move to do and its sequential number in the overall path
	type stackOp struct {
		moveNum int
		pos     position
	}

	var stack = make([]stackOp, 0)

	// Setup initial move
	stack = append(stack, stackOp{moveNum: 0, pos: start})

	for {
		if len(stack) == 0 {
			// No more moves
			break
		}

		// Pop the next move to try
		var op = stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]

		// Handle the backtrack situation when the current branch didn't
		// generate any path and we have to go several moves back
		if op.moveNum < len(s.moves) {
			s.undoMoves(op.moveNum)
		}

		if !s.jump(op.pos) {
			continue
		}

		// We found a path when jumped over every cell on the board
		if len(s.moves) == s.boardSize() {
			return s.moves
		}

		// Otherwise queue potential moves from the current cell
		for _, newPos := range s.getAvailableMoves(op.pos) {
			stack = append(stack, stackOp{moveNum: op.moveNum + 1, pos: newPos})
		}
	}

	return nil
}

// getAvailableMoves returns a list of proper moves which can be made from
// the specified cell.
func (s *solver) getAvailableMoves(p position) []position {
	s.movesBuf = s.movesBuf[:0]
	for _, dir := range s.directions {
		var newPos = p.move(dir)
		if s.canJump(newPos) {
			s.movesBuf = append(s.movesBuf, newPos)
		}
	}

	// Apply Warnsdorff's rule: the move with the least potential next moves
	// should be handled first.
	// Note: since available moves are processed from the highest index to
	// the lowest one they are sorted by descending order of available moves
	// count.
	// var movesCount = make([]int, len(moves))
	for i := range s.movesBuf {
		s.movesCount[i] = s.getPotentialMovesCount(s.movesBuf[i])
	}

	sort.Slice(s.movesBuf, func(i, j int) bool { return s.movesCount[i] > s.movesCount[j] })

	return s.movesBuf
}

// getPotentialMovesCount returns amount of proper moves which can be done
// from the specified cell
func (s *solver) getPotentialMovesCount(p position) int {
	var count = 0
	for _, dir := range s.directions {
		if s.canJump(p.move(dir)) {
			count++
		}
	}
	return count
}

// validateMoveDirection checks that current move is made into the proper
// direction from the prev move
func validateMoveDirection(prev, current position, directions []position) bool {
	for _, dir := range directions {
		if prev.move(dir) == current {
			return true
		}
	}
	return false
}

// validatePath ensures that it forms a proper pawn's path over the whole board
func validatePath(width, height int, directions, moves []position) bool {
	var solver = newSolver(width, height, directions)

	if len(moves) != solver.boardSize() {
		return false
	}

	for i, m := range moves {
		if i > 0 {
			if !validateMoveDirection(moves[i-1], m, directions) {
				return false
			}
		}

		if !solver.jump(m) {
			return false
		}
	}

	return true
}

// Predefined list of possible relative jump directions

var jumpDirections = []position{
	// horizontal
	position{x: 3, y: 0},
	position{x: -3, y: 0},
	// vertical
	position{x: 0, y: 3},
	position{x: 0, y: -3},
	// diagonal
	position{x: 2, y: 2},
	position{x: -2, y: 2},
	position{x: 2, y: -2},
	position{x: -2, y: -2},
}

func main() {
	var width = 10
	var height = 10
	var start = position{x: 1, y: 0}

	var solver = newSolver(width, height, jumpDirections)
	var path = solver.findPath(start)
	if path != nil {
		fmt.Println("Path found")
		for i := range path {
			fmt.Printf("%v\n", path[i])
		}
		var valid = validatePath(width, height, jumpDirections, path)
		if valid {
			fmt.Printf("Path is correct")
		} else {
			fmt.Printf("Path is incorrect")
		}
	} else {
		fmt.Println("Path not found")
	}
}
