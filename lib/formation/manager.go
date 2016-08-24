package formation

import (
	. "github.com/stojg/vivere/lib/vector"
)

type Character interface {
	SetTarget(Static)
	Position() *Vector3
	Orientation() *Quaternion
}

func NewManager(pattern Pattern) *Manager {
	return &Manager{
		slotAssignments: make(SlotAssignments, 0),
		pattern:         pattern,
	}
}

type Manager struct {
	// holds a list of slot assignments
	slotAssignments SlotAssignments
	// holds a Static structure (i.e. Position and Orientation), representing the drift offset of
	// the currently filled slots
	driftOffset Static
	// holds the formation pattern
	pattern Pattern
}

func (m *Manager) AddCharacter(char Character) bool {
	// find out how many slots we have occupied
	occupiedSlots := len(m.slotAssignments)

	if !m.pattern.SupportsSlots(occupiedSlots + 1) {
		return false
	}

	// add new slot assignment
	slotAssignment := &SlotAssignment{
		character: char,
	}
	m.slotAssignments = append(m.slotAssignments, slotAssignment)
	m.updateSlotAssignments()
	return true
}

func (m *Manager) RemoveCharacter(char Character) {
	index, found := m.slotAssignments.find(char)
	if found {
		m.slotAssignments.remove(index)
		m.updateSlotAssignments()
	}
}

func (m *Manager) UpdateSlots() {
	// find the anchor point,
	// @todo, should this be passed in as a closure?
	anchor := m.getAnchorPoint()
	anchorOrientation := anchor.Orientation()

	// go through each character in turn
	for i := range m.slotAssignments {

		// ask for the location of the slot relative to the anchor point, this should be a Static
		relativeLoc := m.pattern.SlotLocation(m.slotAssignments[i].slotNumber)

		// transform it by the anchor points position and orientation
		pos := relativeLoc.Position().Rotate(anchorOrientation).Add(anchor.Position())
		orientation := anchor.Orientation().NewMultiply(relativeLoc.Orientation())

		// remove the drift component
		//pos.Sub(m.driftOffset.Position())
		//orientation.Multiply(m.driftOffset.Orientation().NewInverse())

		m.slotAssignments[i].character.SetTarget(&Model{
			position:    pos,
			orientation: orientation,
		})
	}
}

// updates the assignments of characters to slots
func (m *Manager) updateSlotAssignments() {
	// a very simple assignment algorithm; we simple go through each assignment in the list and
	// assign sequential slot numbers
	for i := range m.slotAssignments {
		m.slotAssignments[i].slotNumber = i
	}
	m.driftOffset = m.pattern.DriftOffset(m.slotAssignments)
}

func (m *Manager) getAnchorPoint() Static {
	anchor := &Model{
		orientation: NewQuaternion(0, 0, 0, 1),
		position:    NewVector3(0, 0, 0),
	}
	for _, assignment := range m.slotAssignments {
		anchor.position.Add(assignment.character.Position())
	}
	anchor.position.Scale(1 / float64(len(m.slotAssignments)))
	return anchor
}
