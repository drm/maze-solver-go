package maze

import "strconv"

func Reporter(state *State, me *Me, m *Maze, viewport Vector) {
	var status string
	if (*state).Outside {
		if !(*state).Solved {
			status = "Looking for entrance ..."
		} else {
			status = "Found it! It took me " + strconv.Itoa(me.NumSteps()) + " steps get there"
		}
	} else {
		status = "Searching ..."
	}
	m.Print(me, viewport, status)
}
