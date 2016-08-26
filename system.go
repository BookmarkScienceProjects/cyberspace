package main

type literal byte

const (
	instEntityID    literal = 1
	instPosition    literal = 2
	instOrientation literal = 3
	instType        literal = 4
	instScale       literal = 5
	instHealth      literal = 6
)

type system interface {
	Update(elapsed float64)
}
