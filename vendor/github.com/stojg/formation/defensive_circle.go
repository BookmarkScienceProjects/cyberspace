package formation

import (
	"github.com/stojg/vector"
	"math"
)

func NewDefensiveCircle(radius float64, initialSlots int) Pattern {
	return &DefensiveCirclePattern{
		numberOfSlots:   initialSlots,
		characterRadius: radius,
	}
}

type DefensiveCirclePattern struct {
	numberOfSlots   int
	characterRadius float64
}

func (p *DefensiveCirclePattern) SetCharacterRadius(r float64) {
	p.characterRadius = r
}

func (p *DefensiveCirclePattern) Radius() float64 {
	if p.characterRadius == 0 || p.numberOfSlots == 0 {
		return 0
	}
	if p.numberOfSlots == 1 {
		return p.characterRadius
	}
	return p.characterRadius / math.Sin(math.Pi/float64(p.numberOfSlots))
}

// Makes sure we can support the given number of slots, in this case we support any number of slots
func (p *DefensiveCirclePattern) SupportsSlots(slotCount int) bool {
	p.numberOfSlots = slotCount
	return true
}

// calculates the drift offset when charactres are in a vicen sets of slots
func (p *DefensiveCirclePattern) DriftOffset(assignments SlotAssignments) Static {

	// store the center of mass
	center := &Model{
		position:    vector.NewVector3(0, 0, 0),
		orientation: vector.NewQuaternion(0, 0, 0, 1),
	}
	// now go through each assignment and add its contribution to the center
	for _, assignment := range assignments {
		location := p.SlotLocation(assignment.slotNumber)
		center.Position().Add(location.Position())
		center.Orientation().Multiply(location.Orientation())
	}
	center.Position().Scale(1 / float64(len(assignments)))
	//center.Orientation().Scale(1/len(assignments))
	return center
}

// Calculates th position of a slot
func (p *DefensiveCirclePattern) SlotLocation(slotNumber int) Static {
	if p.numberOfSlots == 0 {
		panic("number of slots cannot be 0")
	}

	if p.numberOfSlots == 1 {
		return &Model{
			position:    vector.NewVector3(0, 0, 0),
			orientation: vector.NewQuaternion(0, 0, 0, 1),
		}
	}
	// we place the slots around a circle based on their slot number
	angleAroundCircle := float64(slotNumber) / float64(p.numberOfSlots) * math.Pi * 2.0

	// the radios depends on the radius of the character and the number of characters in the circle
	// we want there to be no gap between characters
	radius := p.characterRadius / math.Sin(math.Pi/float64(p.numberOfSlots))
	// create location, and fill its components based on the angle around circle
	location := &Model{
		position: vector.NewVector3(
			radius*math.Cos(angleAroundCircle),
			0,
			radius*math.Sin(angleAroundCircle),
		),
		orientation: vector.NewQuaternion(0, 0, 0, 1),
	}
	return location
}

func (p *DefensiveCirclePattern) calculateNumberOfSlots(assignments SlotAssignments) int {
	// find the number of filled slots: it will be highest slot number in the assignment
	filledSlots := 0
	for _, assignment := range assignments {
		if assignment.slotNumber >= p.numberOfSlots {
			filledSlots = assignment.slotNumber
		}
	}
	// add one to go from the index of the highest slot to the number of the slots needed
	numberOfSlots := filledSlots + 1
	return numberOfSlots
}
