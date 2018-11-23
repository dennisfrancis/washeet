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
// dSrc must be an implementation of SheetDataProvder which is used to draw the contents of the spreadsheet.
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

	if canvasElement == nil || startX+constDefaultCellWidth*10 >= maxX ||
		startY+constDefaultCellHeight*10 >= maxY {
		return nil
	}

	ret := &Sheet{
		document:           js.Global().Get("document"),
		window:             js.Global().Get("window"),
		navigator:          js.Global().Get("navigator"),
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
		paintQueue:         make(chan *sheetPaintRequest, constSheetPaintQueueLength),
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

	// clear the widget area.
	// HACK : maxX + 1.0, maxY + 1.0 is the actual limit
	noStrokeFillRectNoAdjust(&sheet.canvasContext, sheet.origX, sheet.origY, sheet.maxX+1.0, sheet.maxY+1.0, defaultColors.cellFill)
}
