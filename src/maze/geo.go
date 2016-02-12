package maze

// Represents a point in two-dimensional space
type Point struct {
	X, Y int
}

// Represents a relative point
type Vector Point

// Rotate the vector by 90 (1), 180 (0) or 270 (-1) degrees
func (v Vector) turn(dir int) *Vector {
    // if current vertical movement is 0
    if v.Y == 0 { 
		if dir == 0 { // invert horizontal movement
			v.X *= -1
		} else {
			if v.X == 1 { // currently moving west, now move south (or north if dir == -1)
				v.Y = 1 * dir
			} else { // currently moving east, now move north (or south if dir == -1)
				v.Y = -1 * dir 
			}
			v.X = 0
		}
	} else { // current horizontal movement is 0
		if dir == 0 { // invert vertical movement
			v.Y *= -1
		} else {
			if v.Y == 1 { // currently moving south, move west (or east if dir == -1)
				v.X = -1 * dir
			} else { // currently moving north, move east (or west if dir == -1)
				v.X = 1 * dir
			}
			v.Y = 0
		}
	}
	return &v
}

// The hand type is used to identify which hand to keep at the wall
type Hand int

// Returns the "other" hand
func (h Hand) invert() Hand {
    return Hand(int(h) * -1)
}

const (
    LeftHand  Hand = -1
    RightHand Hand = 1
)

