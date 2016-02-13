# Maze solver #
[![Build Status](http://192.168.90.25:8000/api/badges/gerard/dev-challenge-maze-solver/status.svg)](http://192.168.90.25:8000/gerard/dev-challenge-maze-solver)

You can run it by cloning the repo and just

```
export GOPATH=$(pwd)
go run main.go
```

## Flags ##
You can pass the following flags:
```
  -file string
        Load maze from a file in stead of URL
  -fps int
        Frames per second to render (default 6)
  -height int
        Maze height (default 10)
  -left
        Prefer 'left hand'
  -mps int
        Movements per second. Set to 0 to go as fast as possible (default 30)
  -right
        Prefer 'right hand'. Has no effect if left is passed.
  -viewport-height int
        Viewport height (default 40)
  -viewport-width int
        Viewport width (default 80)
  -width int
        Maze width (default 10)
```
