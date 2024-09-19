package blocks

type Inverter struct {
	BaseBlock
	Direction Direction
}

func NewInverter(baseBlock BaseBlock, direction Direction) *Inverter {
	return &Inverter{
		BaseBlock: baseBlock,
		Direction: direction,
	}
}

func (i *Inverter) GetDirection() Direction {
	return i.Direction
}

func (i *Inverter) GetRune(surroundingBlocks map[Direction]Block) rune {
	switch i.Direction {
	case Left:
		return '◀'
	case Right:
		return '▶'
	case Down:
		return '▼'
	case Up:
		return '▲'
	}
	return ' '
}

func (i *Inverter) ConnectsFrom(direction Direction) bool {
	rel := direction

	if (rel == Left || rel == Right) && (i.Direction == Left || i.Direction == Right) {
		return true
	}

	if (rel == Up || rel == Down) && (i.Direction == Up || i.Direction == Down) {
		return true
	}

	return false
}
