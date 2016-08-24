package formation

func NewSlotAssignment(char Character) *SlotAssignment {
	return &SlotAssignment{
		character: char,
	}
}

type SlotAssignment struct {
	character  Character
	slotNumber int
}

func NewSlotAssignments() SlotAssignments {
	return make(SlotAssignments, 0)
}

type SlotAssignments []*SlotAssignment

func (list SlotAssignments) Add(assignment *SlotAssignment) {
	list = append(list, assignment)
}

func (list SlotAssignments) find(char Character) (index int, found bool) {
	for i := range list {
		if list[i].character == char {
			return i, true
		}
	}
	return 0, false
}

func (list SlotAssignments) remove(index int) {
	list = append(list[:index], list[index+1:]...)
}
