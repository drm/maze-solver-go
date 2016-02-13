package maze

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestVectorTurn(t *testing.T) {
    var v, v2, v3 Vector
    v = Vector{0, 1}

    v2 = *v.turn(-1)
    assert.NotEqual(t, v.X, v2.X, "Vectors should differ")
    assert.NotEqual(t, v.Y, v2.Y, "Vectors should differ")
    assert.Equal(t, v2.X, 1, "Direction of X should now be 1")
    assert.Equal(t, v2.Y, 0, "Direction of Y should now be 1")

    v2 = *v.turn(1)
    assert.Equal(t, v2.X, -1, "Direction of X should now be 1")
    assert.Equal(t, v2.Y, 0, "Direction of X should now be 1")

    v2 = *v.turn(0)
    assert.Equal(t, v2.X, 0, "Direction of X should not have changed")
    assert.Equal(t, v2.Y, -1, "Direction of Y should have inverted")

    v.X, v.Y = 1, 0
    v2 = *v.turn(0)
    assert.Equal(t, v2.X, -1, "Direction of X should not have changed")
    assert.Equal(t, v2.Y, 0, "Direction of Y should have inverted")

    v.X, v.Y = 1, 0
    v2 = *v.turn(1).turn(1) 
    v3 = *v.turn(0)
    assert.Equal(t, v2.X, v3.X, "Turning twice should be the same as turning around")
    assert.Equal(t, v2.Y, v3.Y, "Turning twice should be the same as turning around")

    v.X, v.Y = 0, 1
    v2 = *v.turn(1).turn(1) 
    v3 = *v.turn(0)
    assert.Equal(t, v2.X, v3.X, "Turning twice should be the same as turning around")
    assert.Equal(t, v2.Y, v3.Y, "Turning twice should be the same as turning around")
    
    v.X, v.Y = 1, 0
    v2 = *v.turn(-1).turn(-1) 
    v3 = *v.turn(0)
    assert.Equal(t, v2.X, v3.X, "Turning twice should be the same as turning around")
    assert.Equal(t, v2.Y, v3.Y, "Turning twice should be the same as turning around")

    v.X, v.Y = 0, 1
    v2 = *v.turn(-1).turn(-1) 
    v3 = *v.turn(0)
    assert.Equal(t, v2.X, v3.X, "Turning twice should be the same as turning around")
    assert.Equal(t, v2.Y, v3.Y, "Turning twice should be the same as turning around")
}

func TestHand(t *testing.T) {
    l := LeftHand
    assert.Equal(t, RightHand, l.invert(), "Inverted left hand should be equal to RightHand")
    assert.Equal(t, LeftHand, l.invert().invert(), "Inverted left hand twice should be equal to LeftHand")
}
