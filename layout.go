package washeet

import (
	"math"
)

func (self *Sheet) computeLayout() {

	if self == nil {
		return
	}

	// Start of actual top left cell to be drawn after accounting
	// for the row/col-headers
	cellStartX, cellStartY := self.origX+DEFAULT_CELL_WIDTH, self.origY+DEFAULT_CELL_HEIGHT
	currCol, currRow := self.startColumn, self.startRow
	colIdx, rowIdx := 0, 0

	for {

		// We need the startX of the "out-of-screen" column too
		if colIdx >= len(self.colStartXCoords) {
			self.colStartXCoords = append(self.colStartXCoords, cellStartX)
		} else {
			self.colStartXCoords[colIdx] = cellStartX
		}

		if cellStartX > self.maxX {
			break
		}

		cellStartX += math.Max(self.dataSource.GetColumnWidth(currCol), DEFAULT_CELL_WIDTH)
		currCol++
		colIdx++
	}

	for {

		// We need the startY of the "out-of-screen" row too
		if rowIdx >= len(self.rowStartYCoords) {
			self.rowStartYCoords = append(self.rowStartYCoords, cellStartY)
		} else {
			self.rowStartYCoords[rowIdx] = cellStartY
		}

		if cellStartY > self.maxY {
			break
		}

		cellStartY += math.Max(self.dataSource.GetRowHeight(currRow), DEFAULT_CELL_HEIGHT)
		currRow++
		rowIdx++
	}

	self.endColumn, self.endRow = currCol-1, currRow-1
}
