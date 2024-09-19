package blocks

type PoweredBlock struct {
	BaseBlock
}

func NewPoweredBlock(baseBlock BaseBlock) *PoweredBlock {
	return &PoweredBlock{
		BaseBlock: baseBlock,
	}
}

func (pb *PoweredBlock) GetRune(surroundingBlocks map[Direction]Block) rune {
	return '◆'
}
