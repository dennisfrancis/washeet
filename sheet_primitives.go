package washeet

import (
	"math"
)

func (sheet *Sheet) getNearestBorderXY(layout *layoutData, x, y float64, xidx, yidx int64) (bx, by float64, cellxidx, cellyidx int64) {

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

func (sheet *Sheet) getCellIndex(layout *layoutData, x, y float64) (xidx, yidx int64) {

	// last visible column's index
	highx := (layout.endColumn - layout.startColumn)
	// last visible row's index
	highy := (layout.endRow - layout.startRow)

	// sheet.colStartXCoords should have at least highx+2 elements
	// as it stores start coordinate of cell after the last visible.
	xidx = getIntervalIndex(x, layout.colStartXCoords[:highx+2])
	yidx = getIntervalIndex(y, layout.rowStartYCoords[:highy+2])

	return
}

// Warning : no limit checks for args here !
func (sheet *Sheet) getIndices(layout *layoutData, c1, r1, c2, r2 int64) (startXIdx, endXIdx, startYIdx, endYIdx int64) {

	// index of start cell and end cell
	startXIdx = c1 - layout.startColumn
	endXIdx = c2 - layout.startColumn
	// index of start cell and end cell
	startYIdx = r1 - layout.startRow
	endYIdx = r2 - layout.startRow

	return
}

// Warning : no limit checks for args here !
func (sheet *Sheet) getIndicesAndRect(layout *layoutData, c1, r1, c2, r2 int64) (startXIdx, endXIdx, startYIdx, endYIdx int64,
	xlow, xhigh, ylow, yhigh float64) {

	startXIdx, endXIdx, startYIdx, endYIdx = sheet.getIndices(layout, c1, r1, c2, r2)

	xlow = layout.colStartXCoords[startXIdx]
	xhigh = math.Min(layout.colStartXCoords[endXIdx+1], sheet.maxX) // end of last column in view

	ylow = layout.rowStartYCoords[startYIdx]
	yhigh = math.Min(layout.rowStartYCoords[endYIdx+1], sheet.maxY) // end of last row in view

	return
}

func (sheet *Sheet) trimRangeToView(layout *layoutData, colStart int64, rowStart int64, colEnd int64, rowEnd int64) (c1, r1, c2, r2 int64) {

	return maxInt64(colStart, layout.startColumn), maxInt64(rowStart, layout.startRow),
		minInt64(colEnd, layout.endColumn), minInt64(rowEnd, layout.endRow)
}
