package formation

import (
	. "github.com/stojg/vector"
)

// Static could be anything that has a position and Orientation
type Static interface {
	Position() *Vector3
	Orientation() *Quaternion
}

// Character is the minimal interface needed to interact with this package. It represents something
// that has a position, an orientation and can have a target.
type Character interface {
	SetTarget(Static)
	Static
}

// NewManager returns a new formation Manager. It is initialised with anything that can satisfy the
// Pattern interface.
func NewManager(pattern Pattern) *Manager {
	return &Manager{
		slotAssignments: make(SlotAssignments, 0),
		pattern:         pattern,
	}
}

// Manager manages a Pattern. All interactions with a Pattern should be via the Manager. It adds
// and Removes characters and calculates where does characters should be to fit a Patterhappy.
type Manager struct {
	// holds a list of slot assignments
	slotAssignments SlotAssignments
	// holds a Static structure (i.e. Position and Orientation), representing the drift offset of
	// the currently filled slots
	driftOffset Static
	// holds the formation pattern
	pattern Pattern
}

// AddCharacter
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

// RemoveCharacter
func (m *Manager) RemoveCharacter(char Character) {
	index, found := m.slotAssignments.find(char)
	if found {
		m.slotAssignments.remove(index)
		m.updateSlotAssignments()
	}
}

// UpdateSlots
func (m *Manager) UpdateSlots() {
	anchor := m.AnchorPoint()

	anchorOrientation := anchor.Orientation()

	// go through each character in turn
	for _, assignment := range m.slotAssignments {

		// ask for the location of the slot relative to the anchor point, this should be a Static
		relativeLoc := m.pattern.SlotLocation(assignment.slotNumber)

		// transform it by the anchor points position and orientation
		pos := relativeLoc.Position().Rotate(anchorOrientation).Add(anchor.Position())
		orientation := anchor.Orientation().NewMultiply(relativeLoc.Orientation())

		// remove the drift component
		pos.Sub(m.driftOffset.Position())
		//orientation.Multiply(m.driftOffset.Orientation().NewInverse())

		assignment.character.SetTarget(&Model{
			position:    pos,
			orientation: orientation,
		})
	}
}

func (m *Manager) AnchorPoint() Static {
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

// updates the assignments of characters to slots
func (m *Manager) updateSlotAssignments() {
	// a very simple assignment algorithm; we simple go through each assignment in the list and
	// assign sequential slot numbers
	for i := range m.slotAssignments {
		m.slotAssignments[i].slotNumber = i
	}
	m.driftOffset = m.pattern.DriftOffset(m.slotAssignments)
}
