package core

// NewGraphic returns a new Graphic component
func NewGraphic(model int) *Graphic {
	return &Graphic{
		model: model,
	}
}

// Graphic component describes various properties on how it's related GameObject should be rendered
type Graphic struct {
	Component
	isRendered bool
	model      int
}

// IsRendered return true if the related GameObject has been rendered to the screen before
func (g *Graphic) IsRendered() bool {
	return g.isRendered
}

// SetRendered marks the related GameObject as rendered to the screen
func (g *Graphic) SetRendered() {
	g.isRendered = true
}

// Model returns the model number for the related GameObject
func (g *Graphic) Model() int {
	return g.model
}
