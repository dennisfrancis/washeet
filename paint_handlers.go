package washeet

import (
	//	"fmt"
	"math"
)

func (self *Sheet) servePaintRequest(request *sheetPaintRequest) {

	if self == nil || request == nil {
		return
	}

	switch request.kind {
	case sheetPaintWholeSheet:
		self.servePaintWholeSheetRequest(request.col, request.row, request.changeSheetStartCol, request.changeSheetStartRow)
	case sheetPaintCell:
		self.servePaintCellRangeRequest(request.col, request.row, request.col, request.row)
	case sheetPaintCellRange:
		self.servePaintCellRangeRequest(request.col, request.row, request.endCol, request.endRow)
	case sheetPaintSelection:
		self.servePaintSelectionRequest(request.col, request.row, request.endCol, request.endRow)
	}
}

func (self *Sheet) servePaintWholeSheetRequest(col, row int64, changeSheetStartCol, changeSheetStartRow bool) {

	if self == nil {
		return
	}

	rafLD := self.rafLayoutData

	self.computeLayout(rafLD, col, row, changeSheetStartCol, changeSheetStartRow)

	self.drawHeaders(rafLD)

	self.servePaintCellRangeRequest(rafLD.startColumn, rafLD.startRow, rafLD.endColumn, rafLD.endRow)

	self.servePaintSelectionRequest(self.mark.C1, self.mark.R1, self.mark.C2, self.mark.R2)
}

func (self *Sheet) servePaintCellRangeRequest(colStart int64, rowStart int64, colEnd int64, rowEnd int64) {

	if self == nil {
		return
	}

	rafLD := self.rafLayoutData

	c1, r1, c2, r2 := self.trimRangeToView(rafLD, colStart, rowStart, colEnd, rowEnd)

	self.drawRange(rafLD, c1, r1, c2, r2)
}

// Warning : assumes range supplied is well-ordered
func (self *Sheet) servePaintSelectionRequest(colStart, rowStart, colEnd, rowEnd int64) {

	if self == nil {
		return
	}

	rafLD := self.rafLayoutData

	// Undo the current selection
	if !(self.mark.C1 > rafLD.endColumn || self.mark.C2 < rafLD.startColumn || self.mark.R1 > rafLD.endRow || self.mark.R2 < rafLD.startRow) {
		// if current selection is in view at least partially
		c1, r1, c2, r2 := self.trimRangeToView(rafLD, self.mark.C1, self.mark.R1, self.mark.C2, self.mark.R2)
		self.servePaintCellRangeRequest(c1, r1, c2, r2)
	}

	self.mark.C1, self.mark.C2 = colStart, colEnd
	self.mark.R1, self.mark.R2 = rowStart, rowEnd

	// check if mark is out of view
	if self.mark.C1 > rafLD.endColumn || self.mark.C2 < rafLD.startColumn || self.mark.R1 > rafLD.endRow || self.mark.R2 < rafLD.startRow {
		return
	}

	//fmt.Printf("mark = %+v\n", self.mark)

	c1, r1, c2, r2 := self.trimRangeToView(rafLD, self.mark.C1, self.mark.R1, self.mark.C2, self.mark.R2)
	ci1, ci2, ri1, ri2, xlow, xhigh, ylow, yhigh := self.getIndicesAndRect(rafLD, c1, r1, c2, r2)

	if !self.mark.IsSingleCell() {
		strokeFillRect(&self.canvasContext, xlow, ylow, xhigh, yhigh, SELECTION_STROKE_COLOR, SELECTION_FILL_COLOR)
	}

	// Paint borders for the refStartCell if it is in view
	// refStartCell need not be the first cell in the range
	refStartCell := self.selectionState.getRefStartCell()
	if refStartCell.Col <= c2 && refStartCell.Col >= c1 && refStartCell.Row <= r2 && refStartCell.Row >= r1 {

		startCellColIdx, startCellRowIdx := ci1, ri1
		if refStartCell.Col == c2 {
			startCellColIdx = ci2
		}
		if refStartCell.Row == r2 {
			startCellRowIdx = ri2
		}
		xStartCellBeg := math.Max(xlow, rafLD.colStartXCoords[startCellColIdx])
		yStartCellBeg := math.Max(ylow, rafLD.rowStartYCoords[startCellRowIdx])
		xStartCellEnd := math.Min(rafLD.colStartXCoords[startCellColIdx+1], self.maxX)
		yStartCellEnd := math.Min(rafLD.rowStartYCoords[startCellRowIdx+1], self.maxY)
		strokeNoFillRect(&self.canvasContext, xStartCellBeg, yStartCellBeg, xStartCellEnd, yStartCellEnd, CURSOR_STROKE_COLOR)
		strokeNoFillRect(&self.canvasContext, xStartCellBeg+1, yStartCellBeg+1, xStartCellEnd-1, yStartCellEnd-1, CURSOR_STROKE_COLOR)
	}

	if c2 == self.mark.C2 && r2 == self.mark.R2 {
		xLastCellEnd := rafLD.colStartXCoords[ci2+1]
		yLastCellEnd := rafLD.rowStartYCoords[ri2+1]
		if xLastCellEnd <= self.maxX && yLastCellEnd <= self.maxY {
			strokeFillRect(&self.canvasContext, xLastCellEnd-6, yLastCellEnd-6, xLastCellEnd, yLastCellEnd, CURSOR_STROKE_COLOR, CURSOR_STROKE_COLOR)
		}
	}
}
