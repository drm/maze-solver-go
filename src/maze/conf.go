package maze

import (
	"flag"
	"fmt"
)

var Settings Conf

type Conf struct {
	Fps, Mps int
	Viewport Vector
	Maze     Vector
	Left     bool
	Right    bool
	File     string
}

func (conf *Conf) Url() string {
	return fmt.Sprintf("http://www.hereandabove.com/cgi-bin/maze?%d+%d+1+1+5+0+0+0+255+255+255", conf.Maze.X, conf.Maze.Y)
}

func (conf *Conf) Dir() Hand {
	if conf.Left {
		return LeftHand
	}
	return RightHand
}

func init() {
	Settings = Conf{}
	flag.IntVar(&Settings.Fps, "fps", 6, "Frames per second to render")
	flag.IntVar(&Settings.Mps, "mps", 30, "Movements per second. Set to 0 to go as fast as possible")
	flag.IntVar(&Settings.Viewport.X, "viewport-width", 80, "Viewport width")
	flag.IntVar(&Settings.Viewport.Y, "viewport-height", 40, "Viewport height")
	flag.IntVar(&Settings.Maze.X, "width", 10, "Maze width")
	flag.IntVar(&Settings.Maze.Y, "height", 10, "Maze height")
	flag.BoolVar(&Settings.Left, "left", false, "Prefer 'left hand'")
	flag.BoolVar(&Settings.Right, "right", false, "Prefer 'right hand'. Has no effect if left is passed.")
	flag.StringVar(&Settings.File, "file", "", "Load maze from a file in stead of URL")
}

