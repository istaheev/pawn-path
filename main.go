package main

import "fmt"

type position struct {
	x int
	y int
}

type board struct {
	width  int
	height int
}

func (b board) cellIndex(p position) int {
	return p.y*b.width + p.x
}

func (b board) inside(p position) bool {
	return p.x >= 0 && p.x < b.width && p.y >= 0 && p.y < b.height
}

func (b board) size() int {
	return b.width * b.height
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
/*
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
*/
var jumpDirections = []position{
	// horizontal
	position{x: 2, y: 0},
	position{x: -2, y: 0},
	// vertical
	position{x: 0, y: 2},
	position{x: 0, y: -2},
	// diagonal
	position{x: 1, y: 1},
	position{x: -1, y: -1},
}

func findPath(board board, start position) []position {
	// Prepare required data structures
	var visited = make([]bool, board.size())
	var moves = make([]position, board.size())
	var stack = newStack()

	var prevMoveNum = 0
	stack.push(stackOp{moveNum: 0, pos: start})

	for {
		op, ok := stack.pop()
		if !ok {
			// stack is empty
			break
		}

		//fmt.Printf("Move: %d, Prev: %d, Pos: %v\n", op.moveNum, prevMoveNum, op.pos)

		if op.moveNum <= prevMoveNum {
			for i := op.moveNum; i <= prevMoveNum; i++ {
				var p = moves[i]
				visited[board.cellIndex(p)] = false
			}
		}

		var cellIdx = board.cellIndex(op.pos)
		if !board.inside(op.pos) || visited[cellIdx] {
			continue
		}

		visited[cellIdx] = true
		moves[op.moveNum] = op.pos
		prevMoveNum = op.moveNum

		if op.moveNum == board.size()-1 {
			return moves
		}

		for _, dir := range jumpDirections {
			var newPos = position{x: op.pos.x + dir.x, y: op.pos.y + dir.y}
			if !board.inside(newPos) || visited[board.cellIndex(newPos)] {
				continue
			}

			stack.push(stackOp{moveNum: op.moveNum + 1, pos: newPos})
		}
	}

	return nil
}

func validatePath(board board, moves []position) bool {
	// len(moves) == 0 could happen if board.size() == 0
	if len(moves) == 0 || len(moves) != board.size() {
		return false
	}

	var visited = make([]bool, board.size())

	for _, m := range moves {
		if visited[board.cellIndex(m)] {
			return false
		}
		visited[board.cellIndex(m)] = true
	}

	return true
}

func main() {
	var board = board{width: 8, height: 8}
	var start = position{x: 0, y: 0}

	var path = findPath(board, start)
	if path != nil {
		fmt.Println("Path found")
		for i := range path {
			fmt.Printf("%v\n", path[i])
		}
		var valid = validatePath(board, path)
		fmt.Printf("validate: %v", valid)
	} else {
		fmt.Println("Path not found")
	}
}
