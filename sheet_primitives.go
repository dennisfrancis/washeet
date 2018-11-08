package washeet

import (
	"fmt"
	"math"
)

func (self *Sheet) getNearestBorderXY(layout *layoutData, x, y float64, xidx, yidx int64) (bx, by float64, cellxidx, cellyidx int64) {

	bx, by, cellxidx, cellyidx = 0.0, 0.0, -1, -1

	if xidx >= 0 {
		startx := layout.colStartXCoords[xidx]
		endx := layout.colStartXCoords[xidx+1]

		bx = startx
		cellxidx = xidx

		if (endx - x) < (x - startx) {
			bx = endx
			cellxidx = xidx + 1
		}
	}

	if yidx >= 0 {
		starty := layout.rowStartYCoords[yidx]
		endy := layout.rowStartYCoords[yidx+1]

		by = starty
		cellyidx = yidx

		if (endy - y) < (y - starty) {
			by = endy
			cellyidx = yidx + 1
		}
	}

	return
}

func (self *Sheet) getCellIndex(layout *layoutData, x, y float64) (xidx, yidx int64) {

	// last visible column's index
	highx := (layout.endColumn - layout.startColumn)
	// last visible row's index
	highy := (layout.endRow - layout.startRow)

	// self.colStartXCoords should have at least highx+2 elements
	// as it stores start coordinate of cell after the last visible.
	xidx = getIntervalIndex(x, layout.colStartXCoords[:highx+2])
	yidx = getIntervalIndex(y, layout.rowStartYCoords[:highy+2])

	return
}

func (self *Sheet) drawHeaders(layout *layoutData) {

	if self == nil {
		return
	}

	numColsInView := layout.endColumn - layout.startColumn + 1
	numRowsInView := layout.endRow - layout.startRow + 1

	// column header outline
	strokeFillRect(&self.canvasContext, self.origX, self.origY, self.maxX, self.origY+DEFAULT_CELL_HEIGHT, GRID_LINE_COLOR, HEADER_FILL_COLOR)
	// draw column header separators
	drawVertLines(&self.canvasContext, layout.colStartXCoords[0:numColsInView], self.origY, self.origY+DEFAULT_CELL_HEIGHT, GRID_LINE_COLOR)
	// draw col labels (center aligned)
	setFillColor(&self.canvasContext, CELL_DEFAULT_STROKE_COLOR)
	for nCol, nColIdx := layout.startColumn, int64(0); nCol <= layout.endColumn; nCol, nColIdx = nCol+1, nColIdx+1 {
		drawText(&self.canvasContext, layout.colStartXCoords[nColIdx], self.origY,
			layout.colStartXCoords[nColIdx+1], self.origY+DEFAULT_CELL_HEIGHT,
			self.maxX, self.maxY,
			col2ColLabel(nCol), AlignCenter)
	}
	// row header outline
	strokeFillRect(&self.canvasContext, self.origX, self.origY, self.origX+DEFAULT_CELL_WIDTH, self.maxY, GRID_LINE_COLOR, HEADER_FILL_COLOR)
	// draw row header separators
	drawHorizLines(&self.canvasContext, layout.rowStartYCoords[0:numRowsInView], self.origX, self.origX+DEFAULT_CELL_WIDTH, GRID_LINE_COLOR)
	// draw row labels (center aligned)
	setFillColor(&self.canvasContext, CELL_DEFAULT_STROKE_COLOR)
	for nRow, nRowIdx := layout.startRow, int64(0); nRow <= layout.endRow; nRow, nRowIdx = nRow+1, nRowIdx+1 {
		drawText(&self.canvasContext, self.origX, layout.rowStartYCoords[nRowIdx],
			self.origX+DEFAULT_CELL_WIDTH, layout.rowStartYCoords[nRowIdx+1],
			self.maxX, self.maxY,
			row2RowLabel(nRow), AlignCenter)
	}
}

// Warning : no limit check for args here !
func (self *Sheet) drawRange(layout *layoutData, c1, r1, c2, r2 int64) {

	if self == nil {
		return
	}

	startXIdx, endXIdx, startYIdx, endYIdx, xlow, xhigh, ylow, yhigh := self.getIndicesAndRect(layout, c1, r1, c2, r2)

	// cleanup the cell-range area
	noStrokeFillRect(&self.canvasContext, xlow, ylow, xhigh, yhigh, CELL_DEFAULT_FILL_COLOR)

	// draw N vertical lines where N is number of columns in the range
	drawVertLines(&self.canvasContext, layout.colStartXCoords[startXIdx:endXIdx+1], ylow, yhigh, GRID_LINE_COLOR)

	// draw last vertical line to show end of last column
	drawVertLine(&self.canvasContext, ylow, yhigh, xhigh, GRID_LINE_COLOR)

	// draw N horizontal lines where N is number of rows in the range
	drawHorizLines(&self.canvasContext, layout.rowStartYCoords[startYIdx:endYIdx+1], xlow, xhigh, GRID_LINE_COLOR)

	// draw last horizontal line to show end of last row
	drawHorizLine(&self.canvasContext, xlow, xhigh, yhigh, GRID_LINE_COLOR)

	self.drawCellRangeContents(layout, c1, r1, c2, r2)

}

func (self *Sheet) drawCellRangeContents(layout *layoutData, c1, r1, c2, r2 int64) {

	startXIdx, endXIdx, startYIdx, endYIdx := self.getIndices(layout, c1, r1, c2, r2)

	setFillColor(&self.canvasContext, CELL_DEFAULT_STROKE_COLOR)

	for cidx, nCol := startXIdx, c1; cidx <= endXIdx; cidx, nCol = cidx+1, nCol+1 {
		for ridx, nRow := startYIdx, r1; ridx <= endYIdx; ridx, nRow = ridx+1, nRow+1 {

			drawText(&self.canvasContext, layout.colStartXCoords[cidx], layout.rowStartYCoords[ridx],
				layout.colStartXCoords[cidx+1], layout.rowStartYCoords[ridx+1],
				self.maxX, self.maxY,
				self.dataSource.GetDisplayString(nCol, nRow), AlignRight)
		}
	}
}

// Warning : no limit checks for args here !
func (self *Sheet) getIndices(layout *layoutData, c1, r1, c2, r2 int64) (startXIdx, endXIdx, startYIdx, endYIdx int64) {

	// index of start cell and end cell
	startXIdx = c1 - layout.startColumn
	endXIdx = c2 - layout.startColumn
	// index of start cell and end cell
	startYIdx = r1 - layout.startRow
	endYIdx = r2 - layout.startRow

	return
}

// Warning : no limit checks for args here !
func (self *Sheet) getIndicesAndRect(layout *layoutData, c1, r1, c2, r2 int64) (startXIdx, endXIdx, startYIdx, endYIdx int64,
	xlow, xhigh, ylow, yhigh float64) {

	startXIdx, endXIdx, startYIdx, endYIdx = self.getIndices(layout, c1, r1, c2, r2)

	xlow = layout.colStartXCoords[startXIdx]
	xhigh = math.Min(layout.colStartXCoords[endXIdx+1], self.maxX) // end of last column in view

	ylow = layout.rowStartYCoords[startYIdx]
	yhigh = math.Min(layout.rowStartYCoords[endYIdx+1], self.maxY) // end of last row in view

	return
}

func (self *Sheet) trimRangeToView(layout *layoutData, colStart int64, rowStart int64, colEnd int64, rowEnd int64) (c1, r1, c2, r2 int64) {

	return maxInt64(colStart, layout.startColumn), maxInt64(rowStart, layout.startRow),
		minInt64(colEnd, layout.endColumn), minInt64(rowEnd, layout.endRow)
}

func (self *Sheet) addPaintRequest(request *sheetPaintRequest) bool {

	if self == nil || self.stopSignal {
		return false
	}

	queued := false
	select {
	case self.paintQueue <- request:
		queued = true
	default:
		// Queue is full, drop request
		fmt.Printf("[D]")
	}
	return queued
}
