package main

import "fmt"

const W = 64
const H = 64

type Board [W][H]ICell

func (b Board) requestRun(pos Vec2) {
	if !b.contains(pos) {
		return
	}
	v := b.get(pos)
	if v.updated {
		return
	}
	v.updated = true
	v.cell.update(&b, pos)
}
func (b Board) get(pos Vec2) *ICell {
	return &b[pos.x][pos.y]
}

func (b Board) set(pos Vec2, v ICell) {
	b[pos.x][pos.y] = v
}

func (board *Board) tryPushTo(value *int, dir Vec2, pos Vec2) bool {
	sampling := pos.add(dir)
	if !board.contains(sampling) {
		return false
	}
	v := board.get(pos).value
	if v == nil {
		board.get(pos).value = value
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
	return Vec2{x: v1.x + v2.x, y: v2.y + v2.y}
}

type ICell struct {
	value   *int
	cell    Cell
	updated bool
}

type Cell interface {
	symbol() string
	update(board *Board, pos Vec2)
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

func (c Conveyor) update(board *Board, pos Vec2) {
	cur := board.get(pos).value
	if cur == nil {
		return
	}
	success := board.tryPushTo(cur, getDirVec(c.dir), pos)
	if success {
		board.get(pos).value = nil
	}
}

type Generator struct {
	dir   Direction
	value int
}

func (g Generator) symbol() string {
	return "\033[36m" + directionalArrow(g.dir)
}

func (g Generator) update(board *Board, pos Vec2) {
	board.tryPushTo(&g.value, getDirVec(g.dir), pos)
}

type Empty struct{}

func (Empty) symbol() string      { return " " }
func (Empty) update(*Board, Vec2) {}

func main() {
	board := Board{}
	for x := 0; x < W; x++ {
		for y := 0; y < H; y++ {
			board.set(Vec2{x: x, y: y}, ICell{updated: false, value: nil, cell: Empty{}})
		}
	}
	board.set(Vec2{x: 0, y: 0}, ICell{updated: false, value: nil, cell: Generator{dir: Down, value: 1}})
	fmt.Println("test!")
	a := 10
	fmt.Println(a)
}
