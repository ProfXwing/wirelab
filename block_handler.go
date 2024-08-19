package main

import "log"

type BlockType int

const (
	EmptyCursor BlockType = iota
	RedstoneBlock
	RedstoneLamp
	Redstone
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

func (bh *BlockHandler) GetBlockRune(blockType BlockType, x int, y int) rune {
	var isRedstone = func(x, y int) bool {
		return bh.IsBlockType(x, y, Redstone)
	}

	switch blockType {
	case EmptyCursor:
		return 'C'
	case RedstoneBlock:
		return 'R'
	case RedstoneLamp:
		return 'L'
	case Redstone:
		left := isRedstone(x-1, y)
		right := isRedstone(x+1, y)
		down := isRedstone(x, y+1)
		up := isRedstone(x, y-1)

		log.Printf("left: %v\n", left)
		log.Printf("right: %v\n", right)
		log.Printf("up: %v\n", up)
		log.Printf("down: %v\n", left)

		// Four sided
		if left == true && up == true && right == true && down == true {
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
		return bh.IsBlockType(x, y, RedstoneBlock)
	}
	isPoweredRedstone := func(x, y int) bool {
		block := bh.GetBlock(x, y)
		return block != nil && block.BlockType == Redstone && block.Powered
	}

	connectedToRedstoneBlock := testSurroundingBlocks(x, y, isRedstoneBlock)
	connectedToPoweredRedstone := testSurroundingBlocks(x, y, isPoweredRedstone)

	return connectedToRedstoneBlock || connectedToPoweredRedstone
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

func (bh *BlockHandler) NewBlock(cursor *Cursor) {
	x := cursor.X
	y := cursor.Y

	newBlock := &Block{
		BlockType: cursor.SelectedBlockType,
		Powered:   bh.IsBlockPowered(x, y),
		X:         x,
		Y:         y,
	}

	bh.Blocks = append(bh.Blocks, newBlock)
	bh.SetBlock(x, y, newBlock)

	bh.UpdateSurroundingBlocks(x, y)
}
