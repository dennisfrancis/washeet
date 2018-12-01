package washeet

import "fmt"

func (sheet *Sheet) addPaintRequest(request *sheetPaintRequest) bool {

	if sheet == nil || sheet.stopSignal {
		return false
	}

	queued := false
	select {
	case sheet.paintQueue <- request:
		queued = true
	default:
		// Queue is full, drop request
		fmt.Printf("[D]")
	}
	return queued
}

// if col/row = -1 no changes are made before whole-redraw
// changeSheetStartCol/changeSheetStartRow is also used to set sheet.layoutFromStartCol/sheet.layoutFromStartRow
func (sheet *Sheet) paintWholeSheet(col, row int64, changeSheetStartCol, changeSheetStartRow bool) bool {
	req := &sheetPaintRequest{
		kind:                sheetPaintWholeSheet,
		col:                 col,
		row:                 row,
		changeSheetStartCol: changeSheetStartCol,
		changeSheetStartRow: changeSheetStartRow,
	}

	if sheet.addPaintRequest(req) {
		// Layout will be partly/fully computed in RAF thread, so independently compute that here too.
		sheet.computeLayout(sheet.evtHndlrLayoutData, col, row, changeSheetStartCol, changeSheetStartRow)
		return true
	}

	return false
}

func (sheet *Sheet) paintCell(col int64, row int64) bool {

	if sheet == nil {
		return false
	}

	layout := sheet.evtHndlrLayoutData

	// optimization : don't fill the queue with these
	// if we know they are not going to be painted.
	if col < layout.startColumn || col > layout.endColumn ||
		row < layout.startRow || row > layout.endRow {
		return false
	}

	return sheet.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintCell,
		col:    col,
		row:    row,
		endCol: col,
		endRow: row,
	})
}

func (sheet *Sheet) paintCellRange(colStart int64, rowStart int64, colEnd int64, rowEnd int64) bool {

	if sheet == nil {
		return false
	}

	return sheet.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintCellRange,
		col:    colStart,
		row:    rowStart,
		endCol: colEnd,
		endRow: rowEnd,
	})
}

func (sheet *Sheet) paintCellSelection(col, row int64) bool {
	if sheet == nil {
		return false
	}

	return sheet.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintSelection,
		col:    col,
		row:    row,
		endCol: col,
		endRow: row,
	})
}

func (sheet *Sheet) paintCellRangeSelection(colStart, rowStart, colEnd, rowEnd int64) bool {
	if sheet == nil {
		return false
	}

	return sheet.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintSelection,
		col:    colStart,
		row:    rowStart,
		endCol: colEnd,
		endRow: rowEnd,
	})
}

func (sheet *Sheet) drawHeaders(layout *layoutData) {

	if sheet == nil {
		return
	}

	numColsInView := layout.endColumn - layout.startColumn + 1
	numRowsInView := layout.endRow - layout.startRow + 1

	// column header outline
	strokeFillRect(&sheet.canvasContext, sheet.origX, sheet.origY, sheet.maxX,
		sheet.origY+constDefaultCellHeight, defaultColors.gridLine, defaultColors.headerFill)
	// draw column header separators
	drawVertLines(&sheet.canvasContext, layout.colStartXCoords[0:numColsInView], sheet.origY,
		sheet.origY+constDefaultCellHeight, defaultColors.gridLine)
	// draw col labels (center aligned)
	setFillColor(&sheet.canvasContext, defaultColors.cellStroke)
	for nCol, nColIdx := layout.startColumn, int64(0); nCol <= layout.endColumn; nCol, nColIdx = nCol+1, nColIdx+1 {
		drawText(&sheet.canvasContext, layout.colStartXCoords[nColIdx], sheet.origY,
			layout.colStartXCoords[nColIdx+1], sheet.origY+constDefaultCellHeight,
			sheet.maxX, sheet.maxY,
			col2ColLabel(nCol), AlignCenter)
	}
	// row header outline
	strokeFillRect(&sheet.canvasContext, sheet.origX, sheet.origY,
		sheet.origX+constDefaultCellWidth, sheet.maxY, defaultColors.gridLine, defaultColors.headerFill)
	// draw row header separators
	drawHorizLines(&sheet.canvasContext, layout.rowStartYCoords[0:numRowsInView], sheet.origX,
		sheet.origX+constDefaultCellWidth, defaultColors.gridLine)
	// draw row labels (center aligned)
	setFillColor(&sheet.canvasContext, defaultColors.cellStroke)
	for nRow, nRowIdx := layout.startRow, int64(0); nRow <= layout.endRow; nRow, nRowIdx = nRow+1, nRowIdx+1 {
		drawText(&sheet.canvasContext, sheet.origX, layout.rowStartYCoords[nRowIdx],
			sheet.origX+constDefaultCellWidth, layout.rowStartYCoords[nRowIdx+1],
			sheet.maxX, sheet.maxY,
			row2RowLabel(nRow), AlignCenter)
	}
}

// Warning : no limit check for args here !
func (sheet *Sheet) drawRange(layout *layoutData, c1, r1, c2, r2 int64) {

	if sheet == nil {
		return
	}

	startXIdx, endXIdx, startYIdx, endYIdx, xlow, xhigh, ylow, yhigh := sheet.getIndicesAndRect(layout, c1, r1, c2, r2)

	// cleanup the cell-range area
	noStrokeFillRect(&sheet.canvasContext, xlow, ylow, xhigh, yhigh, defaultColors.cellFill)

	// draw N vertical lines where N is number of columns in the range
	drawVertLines(&sheet.canvasContext, layout.colStartXCoords[startXIdx:endXIdx+1], ylow, yhigh, defaultColors.gridLine)

	// draw last vertical line to show end of last column
	drawVertLine(&sheet.canvasContext, ylow, yhigh, xhigh, defaultColors.gridLine)

	// draw N horizontal lines where N is number of rows in the range
	drawHorizLines(&sheet.canvasContext, layout.rowStartYCoords[startYIdx:endYIdx+1], xlow, xhigh, defaultColors.gridLine)

	// draw last horizontal line to show end of last row
	drawHorizLine(&sheet.canvasContext, xlow, xhigh, yhigh, defaultColors.gridLine)

	sheet.drawCellRangeContents(layout, c1, r1, c2, r2)

}

func (sheet *Sheet) drawCellRangeContents(layout *layoutData, c1, r1, c2, r2 int64) {

	startXIdx, endXIdx, startYIdx, endYIdx := sheet.getIndices(layout, c1, r1, c2, r2)

	setFillColor(&sheet.canvasContext, defaultColors.cellStroke)

	for cidx, nCol := startXIdx, c1; cidx <= endXIdx; cidx, nCol = cidx+1, nCol+1 {
		for ridx, nRow := startYIdx, r1; ridx <= endYIdx; ridx, nRow = ridx+1, nRow+1 {

			drawText(&sheet.canvasContext, layout.colStartXCoords[cidx], layout.rowStartYCoords[ridx],
				layout.colStartXCoords[cidx+1], layout.rowStartYCoords[ridx+1],
				sheet.maxX, sheet.maxY,
				sheet.dataSource.GetDisplayString(nCol, nRow), AlignRight)
		}
	}
}
