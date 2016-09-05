package formation

type SlotAssignment struct {
	character  Character
	slotNumber int
}

type SlotAssignments []*SlotAssignment

func (list *SlotAssignments) remove(i int) {
	a := *list
	*list = append(a[:i], a[i+1:]...)
}

func (list SlotAssignments) find(char Character) (index int, found bool) {
	for i := range list {
		if list[i].character == char {
			return i, true
		}
	}
	return 0, false
}
