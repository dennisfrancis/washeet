package washeet

import (
	//	"fmt"
	"math"
)

func (self *Sheet) servePaintRequest(request *SheetPaintRequest) {

	if self == nil || request == nil {
		return
	}

	switch request.Kind {
	case SheetPaintWholeSheet:
		self.servePaintWholeSheetRequest()
	case SheetPaintCell:
		self.servePaintCellRangeRequest(request.Col, request.Row, request.Col, request.Row)
	case SheetPaintCellRange:
		self.servePaintCellRangeRequest(request.Col, request.Row, request.EndCol, request.EndRow)
	case SheetPaintSelection:
		self.servePaintSelectionRequest(request.Col, request.Row, request.EndCol, request.EndRow)
	}
}

func (self *Sheet) servePaintWholeSheetRequest() {

	if self == nil {
		return
	}

	// Recompute endColumn/endRow colStartXCoords/rowStartYCoords
	self.computeLayout()

	self.drawHeaders()

	self.servePaintCellRangeRequest(self.startColumn, self.startRow, self.endColumn, self.endRow)

	self.servePaintSelectionRequest(self.mark.C1, self.mark.R1, self.mark.C2, self.mark.R2)
}

func (self *Sheet) servePaintCellRangeRequest(colStart int64, rowStart int64, colEnd int64, rowEnd int64) {

	if self == nil {
		return
	}

	c1, r1, c2, r2 := self.trimRangeToView(colStart, rowStart, colEnd, rowEnd)

	self.drawRange(c1, r1, c2, r2)
}

// Warning : assumes range supplied is well-ordered
func (self *Sheet) servePaintSelectionRequest(colStart, rowStart, colEnd, rowEnd int64) {

	if self == nil {
		return
	}

	// Undo the current selection
	if !(self.mark.C1 > self.endColumn || self.mark.C2 < self.startColumn || self.mark.R1 > self.endRow || self.mark.R2 < self.startRow) {
		// if current selection is in view at least partially
		c1, r1, c2, r2 := self.trimRangeToView(self.mark.C1, self.mark.R1, self.mark.C2, self.mark.R2)
		self.servePaintCellRangeRequest(c1, r1, c2, r2)
	}

	self.mark.C1, self.mark.C2 = colStart, colEnd
	self.mark.R1, self.mark.R2 = rowStart, rowEnd

	// check if mark is out of view
	if self.mark.C1 > self.endColumn || self.mark.C2 < self.startColumn || self.mark.R1 > self.endRow || self.mark.R2 < self.startRow {
		return
	}

	//fmt.Printf("mark = %+v\n", self.mark)

	c1, r1, c2, r2 := self.trimRangeToView(self.mark.C1, self.mark.R1, self.mark.C2, self.mark.R2)
	ci1, ci2, ri1, ri2, xlow, xhigh, ylow, yhigh := self.getIndicesAndRect(c1, r1, c2, r2)

	if !self.mark.IsSingleCell() {
		strokeFillRect(self.canvasContext, xlow, ylow, xhigh, yhigh, SELECTION_STROKE_COLOR, SELECTION_FILL_COLOR)
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
		xStartCellBeg := math.Max(xlow, self.colStartXCoords[startCellColIdx])
		yStartCellBeg := math.Max(ylow, self.rowStartYCoords[startCellRowIdx])
		xStartCellEnd := math.Min(self.colStartXCoords[startCellColIdx+1], self.maxX)
		yStartCellEnd := math.Min(self.rowStartYCoords[startCellRowIdx+1], self.maxY)
		strokeNoFillRect(self.canvasContext, xStartCellBeg, yStartCellBeg, xStartCellEnd, yStartCellEnd, CURSOR_STROKE_COLOR)
		strokeNoFillRect(self.canvasContext, xStartCellBeg+2, yStartCellBeg+2, xStartCellEnd-2, yStartCellEnd-2, CURSOR_STROKE_COLOR)
	}

	if c2 == self.mark.C2 && r2 == self.mark.R2 {
		xLastCellEnd := self.colStartXCoords[ci2+1]
		yLastCellEnd := self.rowStartYCoords[ri2+1]
		if xLastCellEnd <= self.maxX && yLastCellEnd <= self.maxY {
			strokeFillRect(self.canvasContext, xLastCellEnd-6, yLastCellEnd-6, xLastCellEnd, yLastCellEnd, CURSOR_STROKE_COLOR, CURSOR_STROKE_COLOR)
		}
	}
}
