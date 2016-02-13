package maze

import (
	"errors"
	"time"
)

func Solver(state *State, m *Maze, me *Me, mps int) {
	var dt time.Duration
	if mps > 0 {
		dt = time.Second / time.Duration(mps)
	}

	for !(*state).Solved {
		if moved := me.Move(m); !moved {
			panic(errors.New("I'm stuck! Did the walls move?"))
		}
		if m.Edge(me.Pos) {
			if !(*state).Outside {
				(*state).Solved = true
			}
			(*state).Outside = true
		} else {
			(*state).Outside = false
		}
		if dt > 0 {
			time.Sleep(dt)
		}
	}
}
