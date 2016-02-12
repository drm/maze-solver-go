package maze

import (
	"fmt"
)

type tile int

const (
    TNone     tile = iota
	TWall
	TPath
    TMe
    TVisited
)

var (
    chars map[tile]rune
)

func init() {
    chars = make(map[tile]rune)
    chars[TWall] = '#' 
    chars[TPath] = ' '
    chars[TMe] = '\u263A'
    chars[TNone] = ' '
    chars[TVisited] = '\u00B7'
}

// Maze represents a two-dimensional matrix [][]Px containing the state of each
// [x][y] coordinate, with the dimensions of Dim
type Maze struct {
	Px  [][]tile
	Dim Vector
}

// Zero-construct a new maze of dimensions w x h
func NewMaze(w, h int) *Maze {
	m := Maze{make([][]tile, w), Vector{w, h}}
	for x := 0; x < w; x++ {
		m.Px[x] = make([]tile, h)
	}
	return &m
}

// Get the type of coordinate at (x, y)
func (m *Maze) At(x, y int) (bool, tile) {
	if x < 0 || x >= m.Dim.X || y < 0 || y >= m.Dim.Y {
		return false, TNone
	}
	return true, m.Px[x][y]
}

func (m *Maze) Edge(p Point) bool {
    return p.X == 0 || p.Y == 0 || p.X == m.Dim.X -1 || p.Y == m.Dim.Y -1
}

func (m *Maze) Movable(p Point, v Vector) bool {
    _, t := m.At(p.X + v.X, p.Y + v.Y)
    return t != TWall
}

// return the render character at the position (x, y)
func (m *Maze) ch(x, y int) rune {
	valid, t := m.At(x, y)
	if !valid {
		return chars[TNone]
	} else if t == TWall {
		return chars[TWall]
	} else {
		return chars[TPath]
	}
}

// Print the maze at the given viewport, with the Me object at the center
func (m *Maze) Print(me *Me, viewport Vector, status string) {
	ref := me.Pos
	var ctr = Point{viewport.X / 2, viewport.Y / 2}
	var trans = Vector{ref.X - ctr.X, ref.Y - ctr.Y}

	fmt.Print("\033[H\033[2J")
	for dy := 0; dy <= viewport.Y; dy++ {
		for dx := 0; dx <= viewport.X; dx++ {
			if ctr.X == dx && ctr.Y == dy {
				fmt.Printf("%c", chars[TMe])
			} else {
				x := dx + trans.X
				y := dy + trans.Y
				if me.Visited(x, y) {
					fmt.Printf("%c", chars[TVisited])
				} else {
					fmt.Printf("%c", m.ch(x, y))
				}
			}
		}
		fmt.Print("\n")
	}

	fmt.Println(status)
}

