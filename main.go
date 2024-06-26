package main

import (
	"fmt"
	"strconv"
	"time"
)

const W = 3
const H = 3

type Board [W][H]ICell

func (board *Board) update() {

	for x := 0; x < W; x++ {
		for y := 0; y < W; y++ {
			board[x][y].updated = false
		}
	}
	for x := 0; x < W; x++ {
		for y := 0; y < W; y++ {
			board.requestRun(Vec2{x: x, y: y})
		}
	}
}

func (board *Board) generateImage() string {
	out := ""
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			if board[x][y].value != nil {
				out += strconv.Itoa(*board[x][y].value)
			} else {
				out += board[x][y].cell.symbol()
			}
		}
		out += "\n"
	}
	return out
}

func (board *Board) requestRun(pos Vec2) {
	cell := board.get(pos)
	if cell.updated {
		return
	}
	cell.updated = true
	cell.cell.update(board, pos)
}

func (board *Board) get(pos Vec2) *ICell {
	return &board[pos.x][pos.y]
}

func (board *Board) set(pos Vec2, value ICell) {
	board[pos.x][pos.y] = value
}

// works on invalid pos
func (board *Board) tryPushTo(value *int, shift Vec2, pos Vec2) bool {
	sampling := pos.add(shift)
	if !board.contains(sampling) {
		return false
	}
	board.requestRun(sampling)
	at := board.get(sampling)
	if at.cell.accept(board, sampling) {
		at.value = value
		return true
	}
	return false
}

type Vec2 struct {
	x int
	y int
}

func (Board) contains(v Vec2) bool {
	return v.x >= 0 && v.x < W && v.y >= 0 && v.y < H
}

func (v1 Vec2) add(v2 Vec2) Vec2 {
	return Vec2{x: v1.x + v2.x, y: v1.y + v2.y}
}

type ICell struct {
	value   *int
	cell    Cell
	updated bool
}

type Cell interface {
	symbol() string
	update(board *Board, pos Vec2)
	accept(board *Board, pos Vec2) bool
}
type Direction int8

const (
	Left Direction = iota
	Right
	Up
	Down
)

type Conveyor struct {
	dir Direction
}

func directionalArrow(d Direction) string {
	switch d {
	case Left:
		return "<"
	case Right:
		return ">"
	case Up:
		return "^"
	case Down:
		return "v"
	}
	panic("Invalid direction")
}

func (c Conveyor) symbol() string {
	return directionalArrow(c.dir)
}

func getDirVec(dir Direction) Vec2 {
	switch dir {
	case Left:
		return Vec2{x: -1, y: 0}
	case Right:
		return Vec2{x: 1, y: 0}
	case Up:
		return Vec2{x: 0, y: -1}
	case Down:
		return Vec2{x: 0, y: 1}
	}
	panic("Invalid direction")
}

var c = 0

func (conveyor Conveyor) update(board *Board, pos Vec2) {
	c++
	cur := board.get(pos).value
	if cur == nil {
		return
	}
	success := board.tryPushTo(cur, getDirVec(conveyor.dir), pos)
	if success {
		board.get(pos).value = nil
	}
}

func (Conveyor) accept(board *Board, pos Vec2) bool {
	return board.get(pos).value == nil
}

type Generator struct {
	dir   Direction
	value int
}

func (g Generator) symbol() string {
	return "\033[36m" + directionalArrow(g.dir) + "\033[0m"
}

func (g Generator) update(board *Board, pos Vec2) {
	board.tryPushTo(&g.value, getDirVec(g.dir), pos)
}

func (g Generator) accept(*Board, Vec2) bool { return false }

type Empty struct{}

func (Empty) symbol() string      { return " " }
func (Empty) update(*Board, Vec2) {}
func (Empty) accept(board *Board, pos Vec2) bool {
	return board.get(pos).value == nil
}

func initBoard() Board {
	board := Board{}
	for x := 0; x < W; x++ {
		for y := 0; y < H; y++ {
			board[x][y] = ICell{updated: false, value: nil, cell: Empty{}}
		}
	}
	board[0][0] = ICell{updated: false, value: nil, cell: Generator{dir: Down, value: 1}}
	board[0][1] = ICell{updated: false, value: nil, cell: Conveyor{dir: Right}}
	return board
}

func main() {
	board := initBoard()
	reset := "\033[" + strconv.Itoa(H) + "A\033[" + strconv.Itoa(H) + "D"
	for {
		fmt.Print(board.generateImage())
		fmt.Print(reset)
		board.update()
		time.Sleep(time.Second)
	}
}
