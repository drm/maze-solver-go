package maze

// "Me" is the cute little fellow that walks the maze
// Pos identifies the current position, Dir identifies the intended movement
// vector, and positions is a slice of all previously visited points
type Me struct {
	Pos           Point
	Dir           Vector
	PreferredHand Hand
	positions     []Point
}

// Create a new Me with a preferred direction
func NewMe(dir Hand) *Me {
	var m Me
	if dir == RightHand {
		m = Me{Point{1, 0}, Vector{1, 0}, dir, []Point{}}
	} else {
		m = Me{Point{0, 1}, Vector{0, 1}, dir, []Point{}}
	}
	return &m
}

// Check if the given (x, y) position was previously visited
func (me *Me) Visited(x, y int) bool {
	for _, v := range me.positions {
		if v.X == x && v.Y == y {
			return true
		}
	}
	return false
}

func (me *Me) NumSteps() int {
	return len(me.positions)
}

// Try to turn around. Fails if the tile behind the current position
// is a wall.
func (me *Me) turnAround(maze *Maze) bool {
	return me.try(maze, *me.Dir.turn(0))
}

// Try to turn in the direction dir. Fails if the tile right or left
// (depending the value of `dir`) is a wall
func (me *Me) turn(maze *Maze, dir Hand) bool {
	return me.try(maze, *me.Dir.turn(int(dir)))
}

// Try to move one step in the current direction
// Fails if the next tile in the current direction is a wall
func (me *Me) forward(maze *Maze) bool {
	return me.try(maze, me.Dir)
}

// Do the movement with the specified vector
// If the resultant point is a wall, returns false (no movement occurs)
// If the resultant point is a path, executes the movement and returns true
func (me *Me) try(maze *Maze, dir Vector) bool {
	if maze.Movable(me.Pos, dir) {
		me.Dir = dir
		me.positions = append(me.positions, me.Pos)
		me.Pos.X += dir.X
		me.Pos.Y += dir.Y
		return true
	}
	return false
}

// Move: i.e. do the next step following the algorithm
// of the "Wall follower"
func (me *Me) Move(maze *Maze) bool {
	if me.turn(maze, me.PreferredHand) {
		return true
	}
	if me.forward(maze) {
		return true
	}
	if me.turn(maze, me.PreferredHand.invert()) {
		return true
	}
	if me.turnAround(maze) {
		return true
	}

	return false
}
