package blocks

type Lever struct {
	BaseBlock
}

func NewLever(baseBlock BaseBlock) *Lever {
	return &Lever{
		BaseBlock: baseBlock,
	}
}

func (l *Lever) GetRune(surroundingBlocks map[Direction]Block) rune {
	if l.Powered {
		return '⊓'
	} else {
		return '⊔'
	}
}
