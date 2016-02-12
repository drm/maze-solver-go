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
    "maze"

	_ "image/gif"
)

var conf Conf

type Conf struct {
	Fps, Mps int
	Viewport maze.Vector
	Maze     maze.Vector
	Left     bool
	Right    bool
	File     string
}

func (conf *Conf) Url() string {
	return fmt.Sprintf("http://www.hereandabove.com/cgi-bin/maze?%d+%d+1+1+5+0+0+0+255+255+255", conf.Maze.X, conf.Maze.Y)
}

func (conf *Conf) Dir() maze.Hand {
	if conf.Left {
		return maze.LeftHand
	}
	return maze.RightHand
}

func init() {
	conf = Conf{}
	flag.IntVar(&conf.Fps, "fps", 6, "Frames per second to render")
	flag.IntVar(&conf.Mps, "mps", 30, "Movements per second. Set to 0 to go as fast as possible")
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
	m := maze.NewMaze(extreme.X-origin.X, extreme.Y-origin.Y)

    // initialize the maze based on the colors in the image
	for x := origin.X; x < extreme.X; x++ {
		for y := origin.Y; y < extreme.Y; y++ {
			r, _, _, _ := i.At(x, y).RGBA()
			if r == 0 {
				m.Px[x-origin.X][y-origin.Y] = maze.TWall
			}
		}
	}

    // initialize the funny little guy
	me := maze.NewMe(conf.Dir())

	outside := true
	status := ""
	found := false

	go func() {
		for !found {
			if moved := me.Move(m); !moved {
				panic(errors.New("I'm stuck! Did the walls move?"))
			}
			if m.Edge(me.Pos) {
				if !outside {
					found = true
				}
				outside = true
			} else {
				outside = false
			}
			if conf.Mps > 0 {
				time.Sleep(time.Second / time.Duration(conf.Mps))
			}
		}
	}()

	print := func() {
		if outside {
			if !found {
				status = "Looking for entrance ..."
			} else {
				status = "Found it! It took me " + strconv.Itoa(me.NumSteps()) + " steps get there"
			}
		} else {
			status = "Searching ..."
		}
		m.Print(me, conf.Viewport, status)
	}
	for !found {
		if conf.Fps > 0 {
			go print()
			time.Sleep(time.Second / time.Duration(conf.Fps))
		}
	}
	print()
}
