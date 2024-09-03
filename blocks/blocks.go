package blocks

type BlockType int

const (
	EmptyBlockType BlockType = iota
	PoweredBlockType
	WiredLampType
	WireType
	LeverType
	InverterType
)

type Block interface {
	GetPosition() (int, int)
	GetRune(surroundingBlocks map[Direction]Block) rune
	GetBlockType() BlockType
	ConnectsFrom(direction Direction) bool
	IsPowered() bool
	SetPowered(powered bool)
}

func NewBlock(blockType BlockType, powered bool, x, y int, direction Direction) Block {
	baseBlock := BaseBlock{
		BlockType: blockType,
		X:         x,
		Y:         y,
		Powered:   powered,
	}

	switch blockType {
	case EmptyBlockType:
		return NewEmptyBlock(baseBlock)
	case PoweredBlockType:
		return NewPoweredBlock(baseBlock)
	case WiredLampType:
		return NewWiredLamp(baseBlock)
	case WireType:
		return NewWire(baseBlock)
	case LeverType:
		return NewLever(baseBlock)
	case InverterType:
		return NewInverter(baseBlock, direction)
	}

	return NewEmptyBlock(baseBlock)
}

type BaseBlock struct {
	BlockType BlockType
	X         int
	Y         int
	Powered   bool
}

func (b *BaseBlock) GetPosition() (int, int) {
	return b.X, b.Y
}

func (b *BaseBlock) GetRune(surroundingBlocks map[Direction]Block) rune {
	return ' '
}

func (b *BaseBlock) GetBlockType() BlockType {
	return b.BlockType
}

func (b *BaseBlock) ConnectsFrom(direction Direction) bool {
	return true
}

func (b *BaseBlock) IsPowered() bool {
	return b.Powered
}

func (b *BaseBlock) SetPowered(powered bool) {
	b.Powered = powered
}

func GetRelativeBlockPosition(fromBlock Block, toBlock Block) Direction {
	fromX, fromY := fromBlock.GetPosition()
	toX, toY := toBlock.GetPosition()

	relX := fromX - toX
	relY := fromY - toY

	if relX == -1 && relY == 0 {
		return Left
	}
	if relX == 1 && relY == 0 {
		return Right
	}
	if relX == 0 && relY == -1 {
		return Up
	}
	if relX == 0 && relY == 1 {
		return Down
	}

	return NonSurrounding
}
