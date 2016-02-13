package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"image"
	"net/http"
	"os"
	"time"

	_ "image/gif"

	"maze"
)

var conf maze.Conf

func init() {
    conf = maze.Settings
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

    m := maze.NewMazeFromImage(i)

	// initialize the funny little guy
	me := maze.NewMe(conf.Dir())

	state := maze.State{false, true}

    go maze.Solver(&state, m, me, conf.Mps)
    for !state.Solved {
		if conf.Fps > 0 {
			go maze.Reporter(&state, me, m, conf.Viewport)
			time.Sleep(time.Second / time.Duration(conf.Fps))
		}
    }
	maze.Reporter(&state, me, m, conf.Viewport)
}
