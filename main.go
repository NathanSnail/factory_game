package main

import "fmt"

const W = 64
const H = 64

type Board [W][H]CellPair

type Vec2 struct {
	x int
	y int
}

func boardContains(v Vec2) bool {
	return v.x >= 0 && v.x < W && v.y >= 0 && v.y < H
}

func (v1 Vec2) add(v2 Vec2) Vec2 {
	return Vec2{x: v1.x + v2.x, y: v2.y + v2.y}
}

type CellPair struct {
	value *int
	cell  Cell
}

type Cell interface {
	symbol() string
	update(board *Board, pos Vec2) int
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

func (c Conveyor) symbol() string {
	switch c.dir {
	case Left:
		return "<"
	case Right:
		return ">"
	case Up:
		return "^"
	case Down:
		return "v"
	}
	panic("Invalid conveyor")
}

func tryPushTo(board *Board, value *int, dir Vec2, pos Vec2) bool {
	sampling := pos.add(dir)
	if !boardContains(sampling) {
		return false
	}
	v := board[sampling.x][sampling.y].value
	if v == nil {
		board[sampling.x][sampling.y].value = value
		return true
	}
	return false
}

func (c Conveyor) getDir() Vec2 {
	switch c.dir {
	case Left:
		return Vec2{x: -1, y: 0}
	case Right:
		return Vec2{x: 1, y: 0}
	case Up:
		return Vec2{x: 0, y: -1}
	case Down:
		return Vec2{x: 0, y: 1}
	}
	panic("Invalid conveyor")
}

func (c Conveyor) update(board *Board, pos Vec2) {
	cur := board[pos.x][pos.y].value
	if cur == nil {
		return
	}
	success := tryPushTo(board, cur, c.getDir(), pos)
	if success {
		board[pos.x][pos.y].value = nil
	}
}

func main() {
	fmt.Println("test!")
	a := 10
	fmt.Println(a)
}
