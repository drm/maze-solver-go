package main

import (
	"fmt"
    "image"
    "os"
    "bufio"
    "time"

	_ "image/gif"
)

const (
    WALL = true
    PATH = false

    STR_WALL = "#"
    STR_PATH = " "
    STR_NONE = "~"
)

var (
    VIEWPORT = Vector{30, 30}
)

type Maze struct {
    Px [][]bool
    Dim Rect
}

func NewMaze(w, h int) *Maze {
    m := Maze{make([][]bool, w), Rect{w, h}}
    for x := 0; x < w; x ++ {
        m.Px[x]= make([]bool, h)
    }
    return &m
}

func (m *Maze) At(x, y int) (bool, bool) {
    if x < 0 || x >= m.Dim.W || y < 0 || y >= m.Dim.H {
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

func (m *Maze) Print(ref Point, viewport Vector, status string) {
    var ctr = Point{viewport.X / 2, viewport.Y / 2}
    var trans = Vector{ref.X - ctr.X, ref.Y - ctr.Y}

    fmt.Print("\033[H\033[2J")
    for y := 0; y <= 30; y ++ {
        for x := 0; x <= 30; x ++ {
            if 15 == x && 15 == y {
                fmt.Print("O")
            } else {
                fmt.Print(m.ch(x + trans.X, y + trans.Y))
            }
        }
        fmt.Print("\n")
    }

    fmt.Println(status)
}

type Me struct {
    Pos Point
    Dir Vector
    hand Vector
}

func NewMe() *Me {
    m := Me{Point{1, 0}, Vector{1, 0}, Vector{0, 1}}
    return &m
}

func (me *Me) Move(maze *Maze) {
    if _, t := maze.At(me.Pos.X + me.hand.X, me.Pos.Y + me.hand.Y); t != WALL {
        me.Turn()
    } else {
        retry:
        _, t := maze.At(me.Pos.X + me.Dir.X, me.Pos.Y + me.Dir.Y)

        if t == WALL {
            me.Turn()
            goto retry
        }
    }

    me.Pos.X += me.Dir.X
    me.Pos.Y += me.Dir.Y
}


func (me *Me) Turn() {
    me.Dir = me.hand
    switch {
    case me.hand.Y == 1:
        me.hand.X = -1
        me.hand.Y = 0
    case me.hand.X == -1:
        me.hand.X = 0
        me.hand.Y = -1
    case me.hand.Y == -1:
        me.hand.X = 1
        me.hand.Y = 0
    case me.hand.X == 1:
        me.hand.X = 0
        me.hand.Y = 1
    }
}

type Point struct {
    X, Y int
}
type Vector struct {
    X, Y int
}
type Rect struct {
    W, H int
}

func decode(file string) (image.Image, string, error) {
    f, err := os.Open(file)
    if err != nil {
        return nil, "", err
    }
    defer f.Close()
    
    return image.Decode(bufio.NewReader(f))
}

func main() {
    i, _, err := decode("./maze.gif")
    if nil != err {
        fmt.Println("Error reading file")
    }

    origin, extreme := i.Bounds().Min, i.Bounds().Max
    maze := NewMaze(extreme.X - origin.X, extreme.Y - origin.Y)

    for x := origin.X; x < extreme.X; x++ {
        for y := origin.Y; y < extreme.Y; y ++ {
            r, _, _, _ := i.At(x, y).RGBA()
            if r == 0 {
                maze.Px[x - origin.X][y - origin.Y] = true
            }
        }
    } 
    me := NewMe()
    
    found := false

    for !found {
        time.Sleep(time.Second / 20)
        
        me.Move(maze)
        if me.Pos.X == 0 || me.Pos.Y == 0 || me.Pos.X == maze.Dim.W -1 || me.Pos.Y == maze.Dim.H -1 {
            maze.Print(me.Pos, VIEWPORT, "Looking for entrance ...")
        } else {
            maze.Print(me.Pos, VIEWPORT, "hm .....")
        }
        
    }
}
