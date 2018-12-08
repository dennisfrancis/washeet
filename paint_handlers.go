package washeet

import (
	"math"
)

func (sheet *Sheet) servePaintRequest(request *sheetPaintRequest) {

	if sheet == nil || request == nil {
		return
	}

	switch request.kind {
	case sheetPaintWholeSheet:
		sheet.servePaintWholeSheetRequest(request.col, request.row, request.changeSheetStartCol, request.changeSheetStartRow)
	case sheetPaintCell:
		sheet.servePaintCellRangeRequest(request.col, request.row, request.col, request.row)
	case sheetPaintCellRange:
		sheet.servePaintCellRangeRequest(request.col, request.row, request.endCol, request.endRow)
	case sheetPaintSelection:
		sheet.servePaintSelectionRequest(request.col, request.row, request.endCol, request.endRow)
	}
}

func (sheet *Sheet) servePaintWholeSheetRequest(col, row int64, changeSheetStartCol, changeSheetStartRow bool) {

	if sheet == nil {
		return
	}

	rafLD := sheet.rafLayoutData

	sheet.computeLayout(rafLD, col, row, changeSheetStartCol, changeSheetStartRow)

	sheet.drawHeaders(rafLD)

	sheet.servePaintCellRangeRequest(rafLD.startColumn, rafLD.startRow, rafLD.endColumn, rafLD.endRow)

	sheet.servePaintSelectionRequest(sheet.mark.c1, sheet.mark.r1, sheet.mark.c2, sheet.mark.r2)
}

func (sheet *Sheet) servePaintCellRangeRequest(colStart int64, rowStart int64, colEnd int64, rowEnd int64) {

	if sheet == nil {
		return
	}

	rafLD := sheet.rafLayoutData

	c1, r1, c2, r2 := sheet.trimRangeToView(rafLD, colStart, rowStart, colEnd, rowEnd)

	sheet.drawRange(rafLD, c1, r1, c2, r2)
}

// Warning : assumes range supplied is well-ordered
func (sheet *Sheet) servePaintSelectionRequest(colStart, rowStart, colEnd, rowEnd int64) {

	if sheet == nil {
		return
	}

	rafLD := sheet.rafLayoutData

	// Undo the current selection
	if !(sheet.mark.c1 > rafLD.endColumn || sheet.mark.c2 < rafLD.startColumn || sheet.mark.r1 > rafLD.endRow || sheet.mark.r2 < rafLD.startRow) {
		// if current selection is in view at least partially
		c1, r1, c2, r2 := sheet.trimRangeToView(rafLD, sheet.mark.c1, sheet.mark.r1, sheet.mark.c2, sheet.mark.r2)
		sheet.servePaintCellRangeRequest(c1, r1, c2, r2)
	}

	sheet.mark.c1, sheet.mark.c2 = colStart, colEnd
	sheet.mark.r1, sheet.mark.r2 = rowStart, rowEnd

	// check if mark is out of view
	if sheet.mark.c1 > rafLD.endColumn || sheet.mark.c2 < rafLD.startColumn || sheet.mark.r1 > rafLD.endRow || sheet.mark.r2 < rafLD.startRow {
		return
	}

	//fmt.Printf("mark = %+v\n", sheet.mark)

	c1, r1, c2, r2 := sheet.trimRangeToView(rafLD, sheet.mark.c1, sheet.mark.r1, sheet.mark.c2, sheet.mark.r2)
	ci1, ci2, ri1, ri2, xlow, xhigh, ylow, yhigh := sheet.getIndicesAndRect(rafLD, c1, r1, c2, r2)

	if !sheet.mark.isSingleCell() {
		strokeFillRect(&sheet.canvasContext, xlow, ylow, xhigh, yhigh,
			defaultColors.selectionStroke, defaultColors.selectionFill)
	}

	// Paint borders for the refStartCell if it is in view
	// refStartCell need not be the first cell in the range
	refStartCell := sheet.selectionState.getRefStartCell()
	if refStartCell.col <= c2 && refStartCell.col >= c1 && refStartCell.row <= r2 && refStartCell.row >= r1 {

		startCellColIdx, startCellRowIdx := ci1, ri1
		if refStartCell.col == c2 {
			startCellColIdx = ci2
		}
		if refStartCell.row == r2 {
			startCellRowIdx = ri2
		}
		xStartCellBeg := math.Max(xlow, rafLD.colStartXCoords[startCellColIdx])
		yStartCellBeg := math.Max(ylow, rafLD.rowStartYCoords[startCellRowIdx])
		xStartCellEnd := math.Min(rafLD.colStartXCoords[startCellColIdx+1], sheet.maxX)
		yStartCellEnd := math.Min(rafLD.rowStartYCoords[startCellRowIdx+1], sheet.maxY)
		strokeNoFillRect(&sheet.canvasContext, xStartCellBeg, yStartCellBeg, xStartCellEnd,
			yStartCellEnd, defaultColors.cursorStroke)
		strokeNoFillRect(&sheet.canvasContext, xStartCellBeg+1, yStartCellBeg+1, xStartCellEnd-1,
			yStartCellEnd-1, defaultColors.cursorStroke)
	}

	if c2 == sheet.mark.c2 && r2 == sheet.mark.r2 {
		xLastCellEnd := rafLD.colStartXCoords[ci2+1]
		yLastCellEnd := rafLD.rowStartYCoords[ri2+1]
		if xLastCellEnd <= sheet.maxX && yLastCellEnd <= sheet.maxY {
			strokeFillRect(&sheet.canvasContext, xLastCellEnd-6, yLastCellEnd-6, xLastCellEnd,
				yLastCellEnd, defaultColors.cursorStroke, defaultColors.cursorStroke)
		}
	}
}

// Warning : assumes range supplied is well-ordered
func (sheet *Sheet) servePaintSelectionDiffRequest(colStart, rowStart, colEnd, rowEnd int64) {

	if sheet.mark.c1 == colStart && sheet.mark.r1 == rowStart &&
		sheet.mark.c2 == colEnd && sheet.mark.r2 == rowEnd {
		// No change in selection
		return
	}

	fixedCorner := cornerType(0xff)
	if sheet.mark.c1 == colStart && sheet.mark.r1 == rowStart {
		fixedCorner = cornerTopLeft
	} else if sheet.mark.c2 == colEnd && sheet.mark.r1 == rowStart {
		fixedCorner = cornerTopRight
	} else if sheet.mark.c1 == colStart && sheet.mark.r2 == rowEnd {
		fixedCorner = cornerBottomLeft
	} else if sheet.mark.c2 == colEnd && sheet.mark.r2 == rowEnd {
		fixedCorner = cornerBottomRight
	} else {
		panic("fixed corner finder bug")
	}

	// TODO: implement the selection diff draw logic.
	if fixedCorner == cornerTopLeft {
	}

	// Update mark to new range after everything is finished.
	sheet.mark.c1, sheet.mark.c2 = colStart, colEnd
	sheet.mark.r1, sheet.mark.r2 = rowStart, rowEnd
}
