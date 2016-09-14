package core

func NewGraphic(model int) *Graphic {
	return &Graphic{
		model: model,
	}
}

type Graphic struct {
	Component
	isRendered bool
	model      int
}

func (g *Graphic) IsRendered() bool {
	return g.isRendered
}

func (g *Graphic) SetRendered() {
	g.isRendered = true
}

func (g *Graphic) Model() int {
	return g.model
}
