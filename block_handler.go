package main

type BlockType int

const (
	EmptyCursor BlockType = iota
	PoweredBlock
	WiredLamp
	Wire
	Lever
	Inverter
)

type Direction int

const (
	Left Direction = iota
	Right
	Up
	Down
)

type Block struct {
	BlockType BlockType
	X         int
	Y         int
	Powered   bool
	Direction Direction
}

type BlockHandler struct {
	game   *Game
	board  [][]*Block
	Blocks []*Block
}

func NewBlockHandler(game *Game) *BlockHandler {
	board := make([][]*Block, GameHeight)
	for i := range board {
		board[i] = make([]*Block, GameWidth)
	}

	return &BlockHandler{
		game:  game,
		board: board,
	}
}

func (bh *BlockHandler) IsValidBlockPosition(x, y int) bool {
	if x < 0 || y < 0 || y >= len(bh.board) || x >= len(bh.board[y]) {
		return false
	}

	return true
}

func (bh *BlockHandler) GetBlock(x int, y int) *Block {
	if !bh.IsValidBlockPosition(x, y) {
		return nil
	}

	return bh.board[y][x]
}

func (bh *BlockHandler) SetBlock(x int, y int, block *Block) {
	if !bh.IsValidBlockPosition(x, y) {
		return
	}

	bh.board[y][x] = block
}

func (bh *BlockHandler) IsBlockType(x, y int, blockType BlockType) bool {
	block := bh.GetBlock(x, y)
	return block != nil && block.BlockType == blockType
}

func (bh *BlockHandler) GetBlockRune(block *Block) rune {
	connectableBlockTypes := []BlockType{Wire, WiredLamp, PoweredBlock, Lever}
	var canConnectToBlock = func(x, y int) bool {
		for _, blockType := range connectableBlockTypes {
			if bh.IsBlockType(x, y, blockType) {
				return true
			}
		}
		return false
	}

	x := block.X
	y := block.Y

	switch block.BlockType {
	case EmptyCursor:
		return ' '
	case PoweredBlock:
		return '▲'
	case WiredLamp:
		if block.Powered {
			return '■'
		} else {
			return '□'
		}
	case Lever:
		if block.Powered {
			return '⊓'
		} else {
			return '⊔'
		}
	case Wire:
		{
			left := canConnectToBlock(x-1, y)
			right := canConnectToBlock(x+1, y)
			down := canConnectToBlock(x, y+1)
			up := canConnectToBlock(x, y-1)

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
	case Inverter:
		switch block.Direction {
		case Left:
			return '⊣'
		case Right:
			return '⊢'
		case Down:
			return '⊤'
		case Up:
			return '⊥'
		}
	}

	return ' '
}

func (bh *BlockHandler) IsBlockPowered(x, y int, includeWire bool) bool {
	testSurroundingBlocks := func(x, y int, f func(x, y int) bool) bool {
		return f(x-1, y) || f(x+1, y) || f(x, y-1) || f(x, y+1)
	}
	isPoweredWire := func(x, y int) bool {
		block := bh.GetBlock(x, y)
		return block != nil && block.BlockType == Wire && block.Powered
	}
	isPoweredBlock := func(x, y int) bool {
		return bh.IsBlockType(x, y, PoweredBlock)
	}
	isPoweredLever := func(x, y int) bool {
		block := bh.GetBlock(x, y)
		return block != nil && block.BlockType == Lever && block.Powered
	}

	connectedToPoweredBlock := testSurroundingBlocks(x, y, isPoweredBlock)
	connectedToPoweredLever := testSurroundingBlocks(x, y, isPoweredLever)
	connectedToPoweredWire := testSurroundingBlocks(x, y, isPoweredWire)

	powered := connectedToPoweredBlock || connectedToPoweredLever
	if includeWire {
		powered = powered || connectedToPoweredWire
	}

	return powered
}

func (bh *BlockHandler) UpdateBlock(x, y int) {
	block := bh.GetBlock(x, y)
	if block == nil {
		return
	}

	if block.BlockType == Wire {
		powered, circuit := bh.GetCircuit(block)

		for _, wire := range circuit {
			wire.Powered = powered
			surroundingBlocks := bh.GetSurroundingBlocks(wire.X, wire.Y)

			for _, block := range surroundingBlocks {
				if block.BlockType == WiredLamp {
					block.Powered = bh.IsBlockPowered(block.X, block.Y, true)
				}
			}
		}
	}
}

func (bh *BlockHandler) GetSurroundingBlocks(x, y int) []*Block {
	blocks := []*Block{}

	left := bh.GetBlock(x-1, y)
	right := bh.GetBlock(x+1, y)
	up := bh.GetBlock(x, y-1)
	down := bh.GetBlock(x, y+1)

	for _, block := range []*Block{left, right, up, down} {
		if block != nil {
			blocks = append(blocks, block)
		}
	}

	return blocks
}

func (bh *BlockHandler) UpdateSurroundingBlocks(x, y int) {
	bh.UpdateBlock(x-1, y)
	bh.UpdateBlock(x+1, y)
	bh.UpdateBlock(x, y-1)
	bh.UpdateBlock(x, y+1)
}

func (bh *BlockHandler) GetCircuit(block *Block) (bool, []*Block) {
	wires := []*Block{block}
	powered := false

	for i := 0; i < len(wires); i++ {
		wire := wires[i]

		if bh.IsBlockPowered(wire.X, wire.Y, false) {
			powered = true
		}

		surroundingBlocks := bh.GetSurroundingBlocks(wire.X, wire.Y)

		for _, wire = range surroundingBlocks {
			if wire.BlockType == Wire && !contains(wires, wire) {
				wires = append(wires, wire)
			}
		}
	}

	return powered, wires
}

func (bh *BlockHandler) NewBlock(cursor *Cursor, insertBlock bool) *Block {
	x := cursor.X
	y := cursor.Y

	newBlock := &Block{
		BlockType: cursor.SelectedBlockType,
		Powered:   bh.IsBlockPowered(x, y, true),
		X:         x,
		Y:         y,
		Direction: cursor.Direction,
	}

	if insertBlock {
		bh.Blocks = append(bh.Blocks, newBlock)
		bh.SetBlock(x, y, newBlock)

		bh.UpdateSurroundingBlocks(x, y)
	}

	return newBlock
}
