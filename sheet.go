package washeet

import (
	//	"fmt"
	"math"
	"syscall/js"
)

func NewSheet(canvasElement, context *js.Value, startX float64, startY float64, maxX float64, maxY float64,
	dSrc SheetDataProvider, dSink SheetModelUpdater) *Sheet {

	// HACK : Adjust for line width of 1.0
	maxX -= 1.0
	maxY -= 1.0

	if context == nil || startX+DEFAULT_CELL_WIDTH*10 >= maxX ||
		startY+DEFAULT_CELL_HEIGHT*10 >= maxY {
		return nil
	}

	ret := &Sheet{
		canvasElement:   canvasElement,
		canvasContext:   context,
		origX:           startX,
		origY:           startY,
		maxX:            maxX,
		maxY:            maxY,
		dataSource:      dSrc,
		dataSink:        dSink,
		rafPendingQueue: make(chan js.Value, SHEET_PAINT_QUEUE_LENGTH),
		startColumn:     int64(0),
		startRow:        int64(0),
		endColumn:       int64(0),
		endRow:          int64(0),
		paintQueue:      make(chan *SheetPaintRequest, SHEET_PAINT_QUEUE_LENGTH),
		colStartXCoords: make([]float64, 0, 1+int(math.Ceil((maxX-startX+1)/DEFAULT_CELL_WIDTH))),
		rowStartYCoords: make([]float64, 0, 1+int(math.Ceil((maxY-startY+1)/DEFAULT_CELL_HEIGHT))),
		mark:            MarkData{0, 0, 6, 6},
		stopSignal:      false,
		stopWaitChan:    make(chan bool),
	}

	// TODO : Move these somewhere else
	setFont(ret.canvasContext, "14px serif")
	setLineWidth(ret.canvasContext, 1.0)

	ret.PaintWholeSheet()
	ret.setupClickHandler()
	ret.setupMousemoveHandler()

	return ret
}

func (self *Sheet) Start() {

	if self == nil {
		return
	}

	self.stopSignal = false
	go self.processQueue()
}

func (self *Sheet) Stop() {

	if self == nil || self.stopSignal {
		return
	}
	self.canvasElement.Call("removeEventListener", "click", self.clickHandler)
	self.canvasElement.Call("removeEventListener", "mousemove", self.mousemoveHandler)
	self.canvasElement.Get("style").Set("cursor", "auto")
	self.clickHandler.Release()
	self.mousemoveHandler.Release()
	self.stopSignal = true
	// clear the widget area.
	// HACK : maxX + 1.0, maxY + 1.0 is the actual limit
	noStrokeFillRectNoAdjust(self.canvasContext, self.origX, self.origY, self.maxX+1.0, self.maxY+1.0, CELL_DEFAULT_FILL_COLOR)
	<-self.stopWaitChan
}

func (self *Sheet) PaintWholeSheet() {
	req := &SheetPaintRequest{Kind: SheetPaintWholeSheet}
	self.addPaintRequest(req)
}

func (self *Sheet) PaintCell(col int64, row int64) {

	if self == nil {
		return
	}

	// optimization : don't fill the queue with these
	// if we know they are not going to be painted.
	if col < self.startColumn || col > self.endColumn ||
		row < self.startRow || row > self.endRow {
		return
	}

	self.addPaintRequest(&SheetPaintRequest{
		Kind:   SheetPaintCell,
		Col:    col,
		Row:    row,
		EndCol: col,
		EndRow: row,
	})
}

func (self *Sheet) PaintCellRange(colStart int64, rowStart int64, colEnd int64, rowEnd int64) {

	if self == nil {
		return
	}

	self.addPaintRequest(&SheetPaintRequest{
		Kind:   SheetPaintCellRange,
		Col:    colStart,
		Row:    rowStart,
		EndCol: colEnd,
		EndRow: rowEnd,
	})
}

func (self *Sheet) PaintCellSelection(col, row int64) {
	if self == nil {
		return
	}

	self.mark.C1, self.mark.C2 = col, col
	self.mark.R1, self.mark.R2 = row, row

	self.addPaintRequest(&SheetPaintRequest{Kind: SheetPaintSelection})
}

func (self *Sheet) PaintCellRangeSelection(colStart, rowStart, colEnd, rowEnd int64) {
	if self == nil {
		return
	}

	self.mark.C1, self.mark.C2 = colStart, colEnd
	self.mark.R1, self.mark.R2 = rowStart, rowEnd

	self.addPaintRequest(&SheetPaintRequest{Kind: SheetPaintSelection})
}
