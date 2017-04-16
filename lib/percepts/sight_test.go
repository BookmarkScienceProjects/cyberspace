package percepts

import (
	"github.com/stojg/cyberspace/lib/core"
	"math"
	"testing"
)

func TestDistance(t *testing.T) {
	me := core.NewGameObject("me")
	me.Transform().Position().Set(0, 0, 0)

	friend := core.NewGameObject("friend")
	friend.Transform().Position().Set(10, 0, 0)

	confidence := Distance(me, friend, 20)
	if confidence != 0.5 {
		t.Errorf("I should be able see my friend, confidence %0.2f", confidence)
	}
}

func TestDistanceTooFar(t *testing.T) {
	me := core.NewGameObject("me")
	me.Transform().Position().Set(0, 0, 0)

	friend := core.NewGameObject("friend")
	friend.Transform().Position().Set(10, 0, 0)

	confidence := Distance(me, friend, 5)
	if confidence != 0.0 {
		t.Errorf("I should not be able see my friend, but see her with confidence %0.2f", confidence)
	}
}

func TestInViewConeInFront(t *testing.T) {
	tests := []struct {
		in  [3]float64
		out float64
	}{
		{[3]float64{0, 0, 0}, 1.0},       // on top
		{[3]float64{1, 0, 0}, 1.0},       // in front
		{[3]float64{-1, 0, 0}, 0.0},      // behind
		{[3]float64{1, 0, 1}, 0.5},       // front-left
		{[3]float64{1, 0, -1}, 0.5},      // front-right
		{[3]float64{0.4142, 0, 1}, 0.25}, // front-right
	}

	me := core.NewGameObject("me")
	me.Transform().Position().Set(0, 0, 0)
	for _, test := range tests {
		friend := core.NewGameObject("friend")
		friend.Transform().Position().Set(test.in[0], test.in[1], test.in[2])
		confidence := InViewCone(me, friend, math.Pi)
		if equals(test.out, confidence) {
			t.Errorf("(%1.f,%1.f,%1.f): Expected confidence %0.4f, got %0.4f", test.in[0], test.in[1], test.in[2], test.out, confidence)
		}
	}
}

func equals(a, b float64) bool {
	const epsilon = 0.0001
	return a > (b+epsilon) || a < (b-epsilon)
}
