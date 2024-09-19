package main

import "redstone/blocks"

type BlockHandler struct {
	game   *Game
	board  [][]blocks.Block
	Blocks []blocks.Block
}

func NewBlockHandler(game *Game) *BlockHandler {
	board := make([][]blocks.Block, GameHeight)
	for i := range board {
		board[i] = make([]blocks.Block, GameWidth)
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

func (bh *BlockHandler) GetBlock(x int, y int) blocks.Block {
	if !bh.IsValidBlockPosition(x, y) {
		return nil
	}

	return bh.board[y][x]
}

func (bh *BlockHandler) SetBlock(x int, y int, block blocks.Block) {
	if !bh.IsValidBlockPosition(x, y) {
		return
	}

	bh.board[y][x] = block
}

func (bh *BlockHandler) IsBlockType(x, y int, blockType blocks.BlockType) bool {
	block := bh.GetBlock(x, y)
	return block != nil && block.GetBlockType() == blockType
}

func (bh *BlockHandler) IsBlockPowered(x, y int, includeWire bool) bool {
	testSurroundingBlocks := func(x, y int, f func(x, y int) bool) bool {
		return f(x-1, y) || f(x+1, y) || f(x, y-1) || f(x, y+1)
	}
	isPoweredWire := func(x, y int) bool {
		block := bh.GetBlock(x, y)
		return block != nil && block.GetBlockType() == blocks.WireType && block.IsPowered()
	}
	isPoweredBlock := func(x, y int) bool {
		return bh.IsBlockType(x, y, blocks.PoweredBlockType)
	}
	isPoweredLever := func(x, y int) bool {
		block := bh.GetBlock(x, y)
		return block != nil && block.GetBlockType() == blocks.LeverType && block.IsPowered()
	}
	isConnectedToPoweredInverter := func(x, y int) bool {
		surroundingBlocks := bh.GetSurroundingBlocks(x, y)

		for direction, block := range surroundingBlocks {
			if block == nil {
				continue
			}

			if block.GetBlockType() == blocks.InverterType && block.IsPowered() {
				inverterDirection := block.GetDirection()

				if inverterDirection == blocks.GetOppositeDirection(direction) {
					return true
				}
			}
		}

		return false
	}

	connectedToPoweredBlock := testSurroundingBlocks(x, y, isPoweredBlock)
	connectedToPoweredLever := testSurroundingBlocks(x, y, isPoweredLever)
	connectedToPoweredWire := testSurroundingBlocks(x, y, isPoweredWire)
	connectedToPoweredInverter := isConnectedToPoweredInverter(x, y)

	powered := connectedToPoweredBlock || connectedToPoweredLever || connectedToPoweredInverter
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

	if block.GetBlockType() == blocks.WireType {
		powered, circuit := bh.GetCircuit(block)

		for _, wire := range circuit {
			wire.SetPowered(powered)
			wireX, wireY := wire.GetPosition()
			surroundingBlocks := bh.GetSurroundingBlocks(wireX, wireY)

			for _, block := range surroundingBlocks {
				if block == nil {
					continue
				}

				blockType := block.GetBlockType()
				x, y := block.GetPosition()

				if blockType == blocks.WiredLampType {
					newPowered := bh.IsBlockPowered(x, y, true)
					block.SetPowered(newPowered)
				}

				if blockType == blocks.InverterType {
					inverterSurrounding := bh.GetSurroundingBlocks(x, y)
					direction := block.GetDirection()
					poweredDirection := blocks.GetOppositeDirection(direction)

					if wire == inverterSurrounding[poweredDirection] {
						block.SetPowered(!powered)

						connectedBlock := inverterSurrounding[direction]
						if connectedBlock != nil {
							connectedX, connectedY := inverterSurrounding[direction].GetPosition()
							bh.UpdateBlock(connectedX, connectedY)
						}
					}
				}
			}
		}
	}
}

func (bh *BlockHandler) GetSurroundingBlocks(x, y int) map[blocks.Direction]blocks.Block {
	surroundingBlocks := map[blocks.Direction]blocks.Block{}

	surroundingBlocks[blocks.Left] = bh.GetBlock(x-1, y)
	surroundingBlocks[blocks.Right] = bh.GetBlock(x+1, y)
	surroundingBlocks[blocks.Up] = bh.GetBlock(x, y-1)
	surroundingBlocks[blocks.Down] = bh.GetBlock(x, y+1)

	return surroundingBlocks
}

func (bh *BlockHandler) UpdateSurroundingBlocks(x, y int) {
	bh.UpdateBlock(x-1, y)
	bh.UpdateBlock(x+1, y)
	bh.UpdateBlock(x, y-1)
	bh.UpdateBlock(x, y+1)
}

func (bh *BlockHandler) GetCircuit(block blocks.Block) (bool, []blocks.Block) {
	wires := []blocks.Block{block}
	powered := false

	for i := 0; i < len(wires); i++ {
		wire := wires[i]

		x, y := wire.GetPosition()
		if bh.IsBlockPowered(x, y, false) {
			powered = true
		}

		surroundingBlocks := bh.GetSurroundingBlocks(x, y)

		for _, wire = range surroundingBlocks {
			if wire != nil && wire.GetBlockType() == blocks.WireType && !contains(wires, wire) {
				wires = append(wires, wire)
			}
		}
	}

	return powered, wires
}

func (bh *BlockHandler) NewBlock(cursor *Cursor, insertBlock bool) blocks.Block {
	x := cursor.X
	y := cursor.Y

	blockType := cursor.SelectedBlockType
	powered := bh.IsBlockPowered(x, y, true)
	if blockType == blocks.InverterType {
		powered = !powered
	}

	newBlock := blocks.NewBlock(
		blockType,
		powered,
		x,
		y,
		cursor.Direction,
	)

	if insertBlock {
		bh.Blocks = append(bh.Blocks, newBlock)
		bh.SetBlock(x, y, newBlock)

		bh.UpdateSurroundingBlocks(x, y)
	}

	return newBlock
}
