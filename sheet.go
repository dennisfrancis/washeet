package washeet

import (
	//	"fmt"
	"math"
	"syscall/js"
	"time"
)

const (
	MAXROWCOUNT              int64   = 1048576
	MAXCOLCOUNT              int64   = 16384
	DEFAULT_CELL_WIDTH       float64 = 120.0
	DEFAULT_CELL_HEIGHT      float64 = 30.0
	SHEET_PAINT_QUEUE_LENGTH int     = 10

	GRID_LINE_COLOR           string = "rgba(50, 50, 50, 1.0)"
	CELL_DEFAULT_FILL_COLOR   string = "rgba(255, 255, 255, 1.0)"
	CELL_DEFAULT_STROKE_COLOR string = "rgba(0, 0, 0, 1.0)" // fonts etc.
	CURSOR_STROKE_COLOR       string = "rgba(0, 0, 0, 1.0)"
	HEADER_FILL_COLOR         string = "rgba(200, 200, 200, 1.0)"
	SELECTION_STROKE_COLOR    string = "rgba(0, 0, 50, 0.9)"
	SELECTION_FILL_COLOR      string = "rgba(0, 0, 200, 0.2)"
)

type SheetPaintType byte

const (
	SheetPaintWholeSheet SheetPaintType = iota
	SheetPaintCell
	SheetPaintCellRange
	SheetPaintSelection
)

type SheetPaintRequest struct {
	Kind   SheetPaintType
	Col    int64
	Row    int64
	EndCol int64
	EndRow int64
}

type TextAlignType byte

const (
	AlignLeft TextAlignType = iota
	AlignCenter
	AlignRight
)

type SheetDataProvider interface {
	GetDisplayString(column int64, row int64) string
	GetColumnWidth(column int64) float64
	GetRowHeight(row int64) float64
}

type SheetModelUpdater interface {
	SetColumnWidth(column int64, width float64)
	SetRowHeight(row int64, height float64)
	SetCellContent(row, column int64, content string)
}

type MarkData struct {
	C1 int64
	R1 int64
	C2 int64
	R2 int64
}

func (self *MarkData) IsSingleCell() bool {

	if self == nil {
		return true
	}

	if self.C1 == self.C2 && self.R1 == self.R2 {
		return true
	}

	return false
}

type Sheet struct {
	canvasElement *js.Value
	canvasContext *js.Value
	origX         float64
	origY         float64
	maxX          float64
	maxY          float64

	dataSource SheetDataProvider
	dataSink   SheetModelUpdater

	rafPendingQueue chan js.Value

	startColumn int64
	startRow    int64

	endColumn int64
	endRow    int64

	paintQueue chan *SheetPaintRequest

	colStartXCoords []float64
	rowStartYCoords []float64

	mark         MarkData
	stopSignal   bool
	stopWaitChan chan bool

	clickHandler     js.Callback
	mousemoveHandler js.Callback
}

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

func (self *Sheet) processQueue() {

	if self == nil {
		return
	}

	for {
		select {
		case request := <-self.paintQueue:
			currRFRequest := js.NewCallback(func(args []js.Value) {
				self.servePaintRequest(request)
				<-self.rafPendingQueue
			})
			self.rafPendingQueue <- js.Global().Call("requestAnimationFrame", currRFRequest)
		default:
			if self.stopSignal {
				close(self.rafPendingQueue)
				for reqID := range self.rafPendingQueue {
					js.Global().Call("cancelAnimationFrame", reqID)
				}
				self.stopWaitChan <- true
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

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
		self.servePaintSelectionRequest()
	}
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
