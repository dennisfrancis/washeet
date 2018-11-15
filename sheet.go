package washeet

import (
	//"fmt"
	//"math"
	"syscall/js"
	"time"
)

// NewSheet creates a spreadsheet ui with the given canvas element within a div container
// The new spreadsheet ui will be drawn in the bounding rectangle ((startX, startY), (maxX, maxY))
// where (startX, startY) is the pixel coordinates of top-left corner and
// (maxX, maxY) is the pixel coordinates of bottom-right corner.
// dSrc must be an implementation of SheetDataProvder which is used to draw the contents of the spreadsheet.NewSheet
// dSink must be an implementation of SheetModelUpdater which is used to communicate back the changes in
// contents of the spreadsheet due to user interaction.
//
// Note that the spreadsheet does not become visible after calling NewSheet. For that Start() method
// needs to be called.
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
		rafLayoutData:      newLayoutData(startX, startY, maxX, maxY),
		evtHndlrLayoutData: newLayoutData(startX, startY, maxX, maxY),
		paintQueue:         make(chan *sheetPaintRequest, SHEET_PAINT_QUEUE_LENGTH),
		mark:               markData{0, 0, 0, 0},
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

	// clear the widget area.
	// HACK : maxX + 1.0, maxY + 1.0 is the actual limit
	noStrokeFillRectNoAdjust(&sheet.canvasContext, sheet.origX, sheet.origY, sheet.maxX+1.0, sheet.maxY+1.0, CELL_DEFAULT_FILL_COLOR)
}

// if col/row = -1 no changes are made before whole-redraw
// changeSheetStartCol/changeSheetStartRow is also used to set sheet.layoutFromStartCol/sheet.layoutFromStartRow
func (sheet *Sheet) PaintWholeSheet(col, row int64, changeSheetStartCol, changeSheetStartRow bool) bool {
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

func (sheet *Sheet) PaintCell(col int64, row int64) bool {

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

func (sheet *Sheet) PaintCellRange(colStart int64, rowStart int64, colEnd int64, rowEnd int64) bool {

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

func (sheet *Sheet) PaintCellSelection(col, row int64) bool {
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

func (sheet *Sheet) PaintCellRangeSelection(colStart, rowStart, colEnd, rowEnd int64) bool {
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
