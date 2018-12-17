package washeet

import (
	"syscall/js"
	"time"
)

// NewSheet creates a spreadsheet ui within the given div container
// The new spreadsheet ui will be drawn in a canvas of dimensions (width, height)
// dSrc must be an implementation of SheetDataProvder which is used to draw the contents of the spreadsheet.
//
// Note that the spreadsheet does not become visible after calling NewSheet. For that Start() method
// needs to be called.
func NewSheet(container *js.Value, width float64, height float64, dSrc SheetDataProvider) *Sheet {

	startX, startY := 0.0, 0.0
	maxX, maxY := width, height
	// HACK : Adjust for line width of 1.0
	maxX -= 1.0
	maxY -= 1.0

	if container == nil || startX+constDefaultCellWidth*10 >= maxX ||
		startY+constDefaultCellHeight*10 >= maxY {
		return nil
	}

	ret := &Sheet{
		document:           js.Global().Get("document"),
		window:             js.Global().Get("window"),
		navigator:          js.Global().Get("navigator"),
		container:          container,
		origX:              startX,
		origY:              startY,
		maxX:               maxX,
		maxY:               maxY,
		dataSource:         dSrc,
		rafLayoutData:      newLayoutData(startX, startY, maxX, maxY),
		evtHndlrLayoutData: newLayoutData(startX, startY, maxX, maxY),
		paintQueue:         make(chan *sheetPaintRequest, constSheetPaintQueueLength),
		mark:               markData{0, 0, 0, 0},
		stopSignal:         false,
		stopRequest:        make(chan struct{}),
		mouseState:         defaultMouseState(),
		selectionState:     defaultSelectionState(),
	}

	ret.setupCanvas(width, height)

	// TODO : Move these somewhere else
	setLineWidth(&ret.canvasStore.sheetCanvasContext, 1.0)
	setLineWidth(&ret.canvasStore.selectionCanvasContext, 1.0)

	ret.paintWholeSheet(ret.evtHndlrLayoutData.startColumn, ret.evtHndlrLayoutData.startRow,
		ret.evtHndlrLayoutData.layoutFromStartCol, ret.evtHndlrLayoutData.layoutFromStartRow)
	ret.setupMouseHandlers()
	ret.setupKeyboardHandlers()

	return ret
}

// Start starts the rendering of the spreadsheet and listens to user interactions.
func (sheet *Sheet) Start() {

	if sheet == nil {
		return
	}

	sheet.stopSignal = false
	sheet.launchRenderer()
}

// Stop stops all internal user-event handlers, stops the rendering loop and cleans up
// the portion of canvas designated to Sheet created by NewSheet.
func (sheet *Sheet) Stop() {

	if sheet == nil || sheet.stopSignal {
		return
	}

	sheet.teardownKeyboardHandlers()
	sheet.teardownMouseHandlers()

	sheet.stopSignal = true
	sheet.stopRequest <- struct{}{} // Will block till the rafWorkerCallback senses this
	time.Sleep(100 * time.Millisecond)

	sheet.rafWorkerCallback.Release()
	sheet.teardownCanvas()
}

// UpdateRange redraws the contents of the provided range.
func (sheet *Sheet) UpdateRange(startColumn, startRow, endColumn, endRow int64) {
	sheet.paintCellRange(startColumn, startRow, endColumn, endRow)
}

// UpdateCell redraws the contents of the provided cell location.
func (sheet *Sheet) UpdateCell(column, row int64) {
	sheet.paintCell(column, row)
}
