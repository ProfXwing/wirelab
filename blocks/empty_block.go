package blocks

type EmptyBlock struct {
	BaseBlock
}

func NewEmptyBlock(baseBlock BaseBlock) *EmptyBlock {
	return &EmptyBlock{
		BaseBlock: baseBlock,
	}
}
