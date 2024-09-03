package blocks

type Wire struct {
	BaseBlock
}

func NewWire(baseBlock BaseBlock) *Wire {
	return &Wire{
		BaseBlock: baseBlock,
	}
}

func (w *Wire) GetConnectableBlockTypes() []BlockType {
	return []BlockType{WireType, WiredLampType, PoweredBlockType, LeverType}
}

func (w *Wire) GetRune(surroundingBlocks map[Direction]Block) rune {
	// checks if each of the surrounding blocks connects from the wire's direction
	connects := func(blockDirection Direction, connectsFromDirection Direction) bool {
		block := surroundingBlocks[blockDirection]
		if block == nil {
			return false
		}

		return block.ConnectsFrom(connectsFromDirection)
	}

	left := connects(Left, Right)
	right := connects(Right, Left)
	down := connects(Down, Up)
	up := connects(Up, Down)

	// Four sided
	if left && up && right && down {
		return '╋'
	}

	// Three sided
	if left && up && right {
		return '┻'
	}
	if up && right && down {
		return '┣'
	}
	if right && down && left {
		return '┳'
	}
	if down && left && up {
		return '┫'
	}

	// Two sided, bent
	if up && right {
		return '┗'
	}
	if right && down {
		return '┏'
	}
	if down && left {
		return '┓'
	}
	if left && up {
		return '┛'
	}

	// Two sided, straight
	if up || down {
		return '┃'
	}
	if left || right {
		return '━'
	}

	// Default
	return '╋'
}
