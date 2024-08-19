package main

type BlockType int

const (
	EmptyCursor BlockType = iota
	PoweredBlock
	WiredLamp
	Wire
	Lever
)

type Block struct {
	BlockType BlockType
	X         int
	Y         int
	Powered   bool
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

	return ' '
}

func (bh *BlockHandler) IsBlockPowered(x, y int) bool {
	testSurroundingBlocks := func(x, y int, f func(x, y int) bool) bool {
		return f(x-1, y) || f(x+1, y) || f(x, y-1) || f(x, y+1)
	}
	isRedstoneBlock := func(x, y int) bool {
		return bh.IsBlockType(x, y, PoweredBlock)
	}
	isPoweredRedstone := func(x, y int) bool {
		block := bh.GetBlock(x, y)
		return block != nil && block.BlockType == Wire && block.Powered
	}
	isPoweredLever := func(x, y int) bool {
		block := bh.GetBlock(x, y)
		return block != nil && block.BlockType == Lever && block.Powered
	}

	connectedToRedstoneBlock := testSurroundingBlocks(x, y, isRedstoneBlock)
	connectedToPoweredRedstone := testSurroundingBlocks(x, y, isPoweredRedstone)
	connectedToPoweredLever := testSurroundingBlocks(x, y, isPoweredLever)

	return connectedToRedstoneBlock || connectedToPoweredRedstone || connectedToPoweredLever
}

func (bh *BlockHandler) UpdateBlock(x, y int) {
	block := bh.GetBlock(x, y)
	if block == nil {
		return
	}

	newPowered := bh.IsBlockPowered(x, y)

	if block.Powered != newPowered {
		block.Powered = newPowered
		bh.UpdateSurroundingBlocks(x, y)
	}
}

func (bh *BlockHandler) UpdateSurroundingBlocks(x, y int) {
	bh.UpdateBlock(x-1, y)
	bh.UpdateBlock(x+1, y)
	bh.UpdateBlock(x, y-1)
	bh.UpdateBlock(x, y+1)
}

func (bh *BlockHandler) NewBlock(cursor *Cursor, insertBlock bool) *Block {
	x := cursor.X
	y := cursor.Y

	newBlock := &Block{
		BlockType: cursor.SelectedBlockType,
		Powered:   bh.IsBlockPowered(x, y),
		X:         x,
		Y:         y,
	}

	if insertBlock {
		bh.Blocks = append(bh.Blocks, newBlock)
		bh.SetBlock(x, y, newBlock)

		bh.UpdateSurroundingBlocks(x, y)
	}

	return newBlock
}
