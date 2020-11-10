package main

import "fmt"

type position struct {
	x int
	y int
}

type stackOp struct {
	moveNum int
	pos     position
}

// stack implements simple stack of operations
type stack struct {
	items []stackOp
	count int
}

func newStack() stack {
	return stack{
		items: make([]stackOp, 0),
		count: 0,
	}
}

func (s *stack) push(op stackOp) {
	if s.count == len(s.items) {
		s.items = append(s.items, op)
	} else {
		s.items[s.count] = op
	}

	s.count++
}

func (s *stack) pop() (stackOp, bool) {
	if s.count == 0 {
		return stackOp{}, false
	}

	s.count--
	return s.items[s.count], true
}

// Predefined list of possible relative jump directions

var jumpDirections1 = []position{
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

var jumpDirections2 = []position{
	// horizontal
	position{x: 1, y: 0},
	position{x: -1, y: 0},
	// vertical
	position{x: 0, y: 1},
	position{x: 0, y: -1},
}

type solver struct {
	width   int
	height  int
	visited []bool
	moves   []position
}

func newSolver(width, height int) solver {
	var boardSize = width * height
	return solver{
		width:   width,
		height:  height,
		visited: make([]bool, boardSize),
		moves:   make([]position, 0, boardSize),
	}
}

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

func (s *solver) backtrack(move int) {
	for i := move; i < len(s.moves); i++ {
		var p = s.moves[i]
		s.setVisited(p, false)
	}
	s.moves = s.moves[0:move]
}

func (s *solver) findPath(start position) []position {
	var stack = newStack()

	stack.push(stackOp{moveNum: 0, pos: start})

	for {
		op, ok := stack.pop()
		if !ok {
			// stack is empty
			break
		}

		// fmt.Printf("Move: %d, Prev: %d, Pos: %v\n", op.moveNum, prevMoveNum, op.pos)

		if op.moveNum < len(s.moves) {
			s.backtrack(op.moveNum)
		}

		if !s.jump(op.pos) {
			continue
		}

		if len(s.moves) == s.boardSize() {
			return s.moves
		}

		for _, dir := range jumpDirections1 {
			var newPos = position{x: op.pos.x + dir.x, y: op.pos.y + dir.y}
			if s.canJump(newPos) {
				stack.push(stackOp{moveNum: op.moveNum + 1, pos: newPos})
			}
		}
	}

	return nil
}

func validatePath(width, height int, moves []position) bool {
	var solver = newSolver(width, height)
	// len(moves) == 0 could happen if board.size() == 0
	if len(moves) == 0 || len(moves) != solver.boardSize() {
		return false
	}

	for _, m := range moves {
		if !solver.canJump(m) {
			return false
		}
		solver.setVisited(m, true)
	}

	return true
}

func main() {
	var width = 8
	var height = 8
	var start = position{x: 0, y: 0}

	var solver = newSolver(width, height)
	var path = solver.findPath(start)
	if path != nil {
		fmt.Println("Path found")
		for i := range path {
			fmt.Printf("%v\n", path[i])
		}
		var valid = validatePath(width, height, path)
		fmt.Printf("validate: %v", valid)
	} else {
		fmt.Println("Path not found")
	}
}
