package washeet

import (
	//"fmt"
	//"math"
	"syscall/js"
	"time"
)

func NewSheet(canvasElement, container *js.Value, startX float64, startY float64, maxX float64, maxY float64,
	dSrc SheetDataProvider, dSink SheetModelUpdater) *Sheet {

	// HACK : Adjust for line width of 1.0
	maxX -= 1.0
	maxY -= 1.0

	if canvasElement == nil || startX+DEFAULT_CELL_WIDTH*10 >= maxX ||
		startY+DEFAULT_CELL_HEIGHT*10 >= maxY {
		return nil
	}

	ret := &Sheet{
		document:           js.Global().Get("document"),
		window:             js.Global().Get("window"),
		container:          container,
		canvasElement:      canvasElement,
		canvasContext:      canvasElement.Call("getContext", "2d"),
		origX:              startX,
		origY:              startY,
		maxX:               maxX,
		maxY:               maxY,
		dataSource:         dSrc,
		dataSink:           dSink,
		rafLayoutData:      NewLayoutData(startX, startY, maxX, maxY),
		evtHndlrLayoutData: NewLayoutData(startX, startY, maxX, maxY),
		paintQueue:         make(chan *sheetPaintRequest, SHEET_PAINT_QUEUE_LENGTH),
		mark:               MarkData{0, 0, 0, 0},
		stopSignal:         false,
		stopRequest:        make(chan struct{}),
		mouseState:         defaultMouseState(),
		selectionState:     defaultSelectionState(),
	}

	// TODO : Move these somewhere else
	setFont(&ret.canvasContext, "14px serif")
	setLineWidth(&ret.canvasContext, 1.0)

	ret.setupClipboardTextArea()
	ret.PaintWholeSheet(ret.evtHndlrLayoutData.startColumn, ret.evtHndlrLayoutData.startRow,
		ret.evtHndlrLayoutData.layoutFromStartCol, ret.evtHndlrLayoutData.layoutFromStartRow)
	ret.setupMouseHandlers()
	ret.setupKeyboardHandlers()

	return ret
}

func (self *Sheet) Start() {

	if self == nil {
		return
	}

	self.stopSignal = false
	self.launchRenderer()
}

func (self *Sheet) Stop() {

	if self == nil || self.stopSignal {
		return
	}

	self.teardownKeyboardHandlers()
	self.teardownMouseHandlers()

	self.stopSignal = true
	self.stopRequest <- struct{}{} // Will block till the rafWorkerCallback senses this
	time.Sleep(100 * time.Millisecond)

	self.rafWorkerCallback.Release()

	// clear the widget area.
	// HACK : maxX + 1.0, maxY + 1.0 is the actual limit
	noStrokeFillRectNoAdjust(&self.canvasContext, self.origX, self.origY, self.maxX+1.0, self.maxY+1.0, CELL_DEFAULT_FILL_COLOR)
}

// if col/row = -1 no changes are made before whole-redraw
// changeSheetStartCol/changeSheetStartRow is also used to set self.layoutFromStartCol/self.layoutFromStartRow
func (self *Sheet) PaintWholeSheet(col, row int64, changeSheetStartCol, changeSheetStartRow bool) bool {
	req := &sheetPaintRequest{
		kind:                sheetPaintWholeSheet,
		col:                 col,
		row:                 row,
		changeSheetStartCol: changeSheetStartCol,
		changeSheetStartRow: changeSheetStartRow,
	}

	if self.addPaintRequest(req) {
		// Layout will be partly/fully computed in RAF thread, so independently compute that here too.
		self.computeLayout(self.evtHndlrLayoutData, col, row, changeSheetStartCol, changeSheetStartRow)
		return true
	}

	return false
}

func (self *Sheet) PaintCell(col int64, row int64) bool {

	if self == nil {
		return false
	}

	layout := self.evtHndlrLayoutData

	// optimization : don't fill the queue with these
	// if we know they are not going to be painted.
	if col < layout.startColumn || col > layout.endColumn ||
		row < layout.startRow || row > layout.endRow {
		return false
	}

	return self.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintCell,
		col:    col,
		row:    row,
		endCol: col,
		endRow: row,
	})
}

func (self *Sheet) PaintCellRange(colStart int64, rowStart int64, colEnd int64, rowEnd int64) bool {

	if self == nil {
		return false
	}

	return self.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintCellRange,
		col:    colStart,
		row:    rowStart,
		endCol: colEnd,
		endRow: rowEnd,
	})
}

func (self *Sheet) PaintCellSelection(col, row int64) bool {
	if self == nil {
		return false
	}

	return self.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintSelection,
		col:    col,
		row:    row,
		endCol: col,
		endRow: row,
	})
}

func (self *Sheet) PaintCellRangeSelection(colStart, rowStart, colEnd, rowEnd int64) bool {
	if self == nil {
		return false
	}

	return self.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintSelection,
		col:    colStart,
		row:    rowStart,
		endCol: colEnd,
		endRow: rowEnd,
	})
}
