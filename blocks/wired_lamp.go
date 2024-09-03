package blocks

type WiredLamp struct {
	BaseBlock
}

func NewWiredLamp(baseBlock BaseBlock) *WiredLamp {
	return &WiredLamp{
		BaseBlock: baseBlock,
	}
}

func (wl *WiredLamp) GetRune(surroundingBlocks map[Direction]Block) rune {
	if wl.Powered {
		return '■'
	} else {
		return '□'
	}
}
