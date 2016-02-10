package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"image"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "image/gif"
)

const (
	WALL = true
	PATH = false

	STR_ME      = "\u263A"
	STR_WALL    = "#"
	STR_PATH    = " "
	STR_NONE    = " "
	STR_VISITED = "\u00B7"
)

var (
	conf Conf
)

type Maze struct {
	Px  [][]bool
	Dim Vector
}

func NewMaze(w, h int) *Maze {
	m := Maze{make([][]bool, w), Vector{w, h}}
	for x := 0; x < w; x++ {
		m.Px[x] = make([]bool, h)
	}
	return &m
}

func (m *Maze) At(x, y int) (bool, bool) {
	if x < 0 || x >= m.Dim.X || y < 0 || y >= m.Dim.Y {
		return false, false
	}
	return true, m.Px[x][y]
}

func (m *Maze) ch(x, y int) string {
	valid, t := m.At(x, y)
	if !valid {
		return STR_NONE
	} else if t == WALL {
		return STR_WALL
	} else {
		return STR_PATH
	}
}

func (m *Maze) Print(me *Me, viewport Vector, status string) {
	ref := me.Pos
	var ctr = Point{viewport.X / 2, viewport.Y / 2}
	var trans = Vector{ref.X - ctr.X, ref.Y - ctr.Y}

	fmt.Print("\033[H\033[2J")
	for dy := 0; dy <= viewport.Y; dy++ {
		for dx := 0; dx <= viewport.X; dx++ {
			if ctr.X == dx && ctr.Y == dy {
				fmt.Print(STR_ME)
			} else {
				x := dx + trans.X
				y := dy + trans.Y
				if me.Visited(x, y) {
					fmt.Print(STR_VISITED)
				} else {
					fmt.Print(m.ch(x, y))
				}
			}
		}
		fmt.Print("\n")
	}

	fmt.Println(status)
}

type Me struct {
	Pos       Point
	Dir       Vector
	positions []Point
}

func NewMe(dir int) *Me {
	var m Me
	if dir == 1 {
		m = Me{Point{1, 0}, Vector{1, 0}, []Point{}}
	} else {
		m = Me{Point{0, 1}, Vector{0, 1}, []Point{}}
	}
	return &m
}

func (me *Me) Visited(x, y int) bool {
	for _, v := range me.positions {
		if v.X == x && v.Y == y {
			return true
		}
	}
	return false
}

func (me *Me) turnAround(maze *Maze) bool {
	return me.try(maze, *me.Dir.turn(0))
}

func (me *Me) turn(maze *Maze, dir int) bool {
	return me.try(maze, *me.Dir.turn(dir))
}

func (me *Me) forward(maze *Maze) bool {
	return me.try(maze, me.Dir)
}

func (me *Me) try(maze *Maze, dir Vector) bool {
	if _, t := maze.At(me.Pos.X+dir.X, me.Pos.Y+dir.Y); t != WALL {
		me.Dir = dir
		me.positions = append(me.positions, me.Pos)
		me.Pos.X += dir.X
		me.Pos.Y += dir.Y
		return true
	}
	return false
}

func (me *Me) Move(maze *Maze) bool {
	if me.turn(maze, conf.Dir()) {
		return true
	}
	if me.forward(maze) {
		return true
	}
	if me.turn(maze, conf.Dir()*-1) {
		return true
	}
	if me.turnAround(maze) {
		return true
	}

	return false
}

type Point struct {
	X, Y int
}
type Vector struct {
	X, Y int
}

func (v Vector) turn(dir int) *Vector {
	if v.X != 0 {
		if dir == 0 {
			v.X *= -1
		} else {
			if v.X == 1 {
				v.Y = 1 * dir
			} else {
				v.Y = -1 * dir
			}
			v.X = 0
		}
	} else {
		if dir == 0 {
			v.Y *= -1
		} else {
			if v.Y == 1 {
				v.X = -1 * dir
			} else {
				v.X = 1 * dir
			}
			v.Y = 0
		}
	}
	return &v
}

type Conf struct {
	Fps, Mps time.Duration
	Viewport Vector
	Maze     Vector
	Left     bool
	Right    bool
	File     string
}

func (conf *Conf) Url() string {
	return fmt.Sprintf("http://www.hereandabove.com/cgi-bin/maze?%d+%d+1+1+5+0+0+0+255+255+255", conf.Maze.X, conf.Maze.Y)
}

func (conf *Conf) Dir() int {
	if conf.Left {
		return -1
	}
	return 1
}

func init() {
	conf = Conf{}
	flag.DurationVar(&conf.Fps, "fps", 6, "Frames per second to render")
	flag.DurationVar(&conf.Mps, "mps", 30, "Movements per second. Set to 0 to go as fast as possible")
	flag.IntVar(&conf.Viewport.X, "viewport-width", 80, "Viewport width")
	flag.IntVar(&conf.Viewport.Y, "viewport-height", 40, "Viewport height")
	flag.IntVar(&conf.Maze.X, "width", 10, "Maze width")
	flag.IntVar(&conf.Maze.Y, "height", 10, "Maze height")
	flag.BoolVar(&conf.Left, "left", false, "Prefer 'left hand'")
	flag.BoolVar(&conf.Right, "right", false, "Prefer 'right hand'. Has no effect if left is passed.")
	flag.StringVar(&conf.File, "file", "", "Load maze from a file in stead of URL")
}

func main() {
	var reader func() (image.Image, string, error)
	flag.Parse()

	if conf.File != "" {
		reader = func() (image.Image, string, error) {
			f, err := os.Open(conf.File)
			if err != nil {
				panic(errors.New("Could not read file " + conf.File))
			}
			defer f.Close()

			return image.Decode(bufio.NewReader(f))
		}
	} else {
		reader = func() (image.Image, string, error) {
			response, err := http.Get(conf.Url())

			if err != nil {
				panic(errors.New("Could not read url " + conf.Url()))
			}

			return image.Decode(response.Body)
		}
	}

	i, _, err := reader()
	if nil != err {
		fmt.Println("Error reading file", err)
		return
	}

	origin, extreme := i.Bounds().Min, i.Bounds().Max
	maze := NewMaze(extreme.X-origin.X, extreme.Y-origin.Y)
	for x := origin.X; x < extreme.X; x++ {
		for y := origin.Y; y < extreme.Y; y++ {
			r, _, _, _ := i.At(x, y).RGBA()
			if r == 0 {
				maze.Px[x-origin.X][y-origin.Y] = true
			}
		}
	}
	me := NewMe(conf.Dir())

	outside := true
	status := ""
	found := false

	go func() {
		for !found {
			if moved := me.Move(maze); !moved {
				panic(errors.New("I'm stuck! Did the walls move?"))
			}
			if me.Pos.X == 0 || me.Pos.Y == 0 || me.Pos.X == maze.Dim.X-1 || me.Pos.Y == maze.Dim.Y-1 {
				if !outside {
					found = true
				}
				outside = true
			} else {
				outside = false
			}
			if conf.Mps > 0 {
				time.Sleep(time.Second / conf.Mps)
			}
		}
	}()

	print := func() {
		if outside {
			if !found {
				status = "Looking for entrance ..."
			} else {
				status = "Found it! It took me " + strconv.Itoa(len(me.positions)) + " steps get there"
			}
		} else {
			status = "Searching ..."
		}
		maze.Print(me, conf.Viewport, status)
	}
	for !found {
		if conf.Fps > 0 {
			go print()
			time.Sleep(time.Second / conf.Fps)
		}
	}
	print()
}
