package core

import "github.com/stojg/goap"

type AI interface {
	SetWorldState(goap.StateList)
	SetGameObject(*GameObject)
	Update(float64)
}
