package washeet

import (
	"math"
)

func NewLayoutData(originX, originY, maxX, maxY float64) *layoutData {

	if maxX <= originX || maxY <= originY {
		return nil
	}

	return &layoutData{
		startColumn:     int64(0),
		startRow:        int64(0),
		endColumn:       int64(0),
		endRow:          int64(0),
		colStartXCoords: make([]float64, 0, 1+int(math.Ceil((maxX-originX+1)/DEFAULT_CELL_WIDTH))),
		rowStartYCoords: make([]float64, 0, 1+int(math.Ceil((maxY-originY+1)/DEFAULT_CELL_HEIGHT))),
	}
}

func (self *Sheet) computeLayout() {

	if self == nil {
		return
	}

	// Start of actual top left cell to be drawn after accounting
	// for the row/col-headers
	minX, minY := self.origX+DEFAULT_CELL_WIDTH, self.origY+DEFAULT_CELL_HEIGHT

	if self.layoutFromStartCol {
		self.endColumn, self.colStartXCoords = computeCellsCoordsRefStart(
			minX,
			self.maxX,
			self.startColumn,
			self.dataSource.GetColumnWidth,
			DEFAULT_CELL_WIDTH,
			self.colStartXCoords,
		)
	} else {
		self.startColumn, self.endColumn, self.colStartXCoords = computeCellsCoordsRefEnd(
			minX,
			self.maxX,
			self.endColumn,
			self.dataSource.GetColumnWidth,
			DEFAULT_CELL_WIDTH,
			self.colStartXCoords,
		)
	}

	if self.layoutFromStartRow {
		self.endRow, self.rowStartYCoords = computeCellsCoordsRefStart(
			minY,
			self.maxY,
			self.startRow,
			self.dataSource.GetRowHeight,
			DEFAULT_CELL_HEIGHT,
			self.rowStartYCoords,
		)
	} else {
		self.startRow, self.endRow, self.rowStartYCoords = computeCellsCoordsRefEnd(
			minY,
			self.maxY,
			self.endRow,
			self.dataSource.GetRowHeight,
			DEFAULT_CELL_HEIGHT,
			self.rowStartYCoords,
		)
	}
}
