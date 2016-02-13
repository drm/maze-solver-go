package maze

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestVectorTurn(t *testing.T) {
    v := Vector{0, 1}

    v2 := v.turn(-1)

    assert.Equal(t, v2.X, 1, "Direction of X should now be 1")
    assert.Equal(t, v2.Y, 0, "Direction of X should now be 1")

    v2 := v.turn(1)

    assert.Equal(t, v2.X, -1, "Direction of X should now be 1")
    assert.Equal(t, v2.Y, 0, "Direction of X should now be 1")
}
