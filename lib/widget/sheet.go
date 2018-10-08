package widget

import (
	"math"
	"syscall/js"
)

const (
	MAXROWCOUNT              int64   = 1048576
	MAXCOLCOUNT              int64   = 16384
	DEFAULT_CELL_WIDTH       float64 = 120.0
	DEFAULT_CELL_HEIGHT      float64 = 50.0
	SHEET_PAINT_QUEUE_LENGTH int     = 10

	GRID_LINE_COLOR           string = "rgba(50, 50, 50, 1.0)"
	CELL_DEFAULT_FILL_COLOR   string = "rgba(255, 255, 255, 1.0)"
	CELL_DEFAULT_STROKE_COLOR string = "rgba(0, 0, 0, 1.0)" // fonts etc.
	CURSOR_STROKE_COLOR       string = "rgba(0, 0, 0, 1.0)"
	HEADER_FILL_COLOR         string = "rgba(100, 100, 100, 1.0)"
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

type SheetDataProvider interface {
	GetDisplayString(column int64, row int64) string
	GetColumnWidth(column int64) float64
	GetRowHeight(row int64) float64
}

type SheetModelUpdater interface {
	SetColumnWidth(column int64, width float64)
	SetRowHeight(row int64, height float64)
	SetCellContent(content string)
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
	canvasContext *js.Value
	origX         float64
	origY         float64
	maxX          float64
	maxY          float64

	dataSource SheetDataProvider
	dataSink   SheetModelUpdater

	renderFrame *js.Callback

	startColumn int64
	startRow    int64

	endColumn int64
	endRow    int64

	paintQueue chan *SheetPaintRequest

	colStartXCoords []float64
	rowStartYCoords []float64

	mark MarkData
}

func NewSheet(context *js.Value, startX float64, startY float64, maxX float64, maxY float64,
	dSrc SheetDataProvider, dSink SheetModelUpdater) *Sheet {

	if context == nil || startX+DEFAULT_CELL_WIDTH*10 >= maxX ||
		startY+DEFAULT_CELL_HEIGHT*10 >= maxY {
		return nil
	}

	ret := &Sheet{
		canvasContext:   context,
		origX:           startX,
		origY:           startY,
		maxX:            maxX,
		maxY:            maxY,
		dataSource:      dSrc,
		dataSink:        dSink,
		renderFrame:     nil,
		startColumn:     int64(0),
		startRow:        int64(0),
		endColumn:       int64(0),
		endRow:          int64(0),
		paintQueue:      make(chan *SheetPaintRequest, SHEET_PAINT_QUEUE_LENGTH),
		colStartXCoords: make([]float64, 0, ceil_int64(maxX-startX+1, DEFAULT_CELL_WIDTH)),
		rowStartYCoords: make([]float64, 0, ceil_int64(maxY-startY+1, DEFAULT_CELL_HEIGHT)),
		mark:            MarkData{0, 0, 0, 0},
	}

	// Compute endColumn/endRow colStartXCoords/rowStartYCoords
	ret.computeLayout()

	return ret
}

func (self *Sheet) computeLayout() {

	if self == nil {
		return
	}

	// Start of actual top left cell to be drawn after accounting
	// for the row/col-headers
	cellStartX, cellStartY := self.origX+DEFAULT_CELL_WIDTH, self.origY+DEFAULT_CELL_HEIGHT
	currCol, currRow := self.startColumn, self.startRow
	colIdx, rowIdx := 0, 0

	for {

		// We need the startX of the "out-of-screen" column too
		if colIdx >= len(self.colStartXCoords) {
			self.colStartXCoords = append(self.colStartXCoords, cellStartX)
		} else {
			self.colStartXCoords[colIdx] = cellStartX
		}

		if cellStartX > self.maxX {
			break
		}

		cellStartX += self.dataSource.GetColumnWidth(currCol)
		currCol++
		colIdx++
	}

	for {

		// We need the startY of the "out-of-screen" row too
		if rowIdx >= len(self.rowStartYCoords) {
			self.rowStartYCoords = append(self.rowStartYCoords, cellStartY)
		} else {
			self.rowStartYCoords[rowIdx] = cellStartY
		}

		if cellStartY > self.maxY {
			break
		}

		cellStartY += self.dataSource.GetRowHeight(currRow)
		currRow++
		rowIdx++
	}

	self.endColumn, self.endRow = currCol-1, currRow-1
}

func (self *Sheet) Start() {

	if self == nil || self.renderFrame != nil {
		return
	}

	self.renderFrame = &js.NewCallback(func(args []js.Value) {
		self.processQueue()
		js.Global().Call("requestAnimationFrame", *self.renderFrame)
	})

	js.Global().Call("requestAnimationFrame", *self.renderFrame)
}

func (self *Sheet) Stop() {

	if self == nil || self.renderFrame == nil {
		return
	}
	self.renderFrame.Release()
	self.renderFrame = nil
	// TODO : clear the widget area.
}

func (self *Sheet) processQueue() {

	if self == nil {
		return
	}

	var request *SheetPaintRequest = nil

	// Process all requests on the queue and return
	for {
		select {
		case request <- self.paintQueue:
			self.servePaintRequest(request)
		default:
			return
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

func (self *Sheet) servePaintWholeSheetRequest() {

	if self == nil {
		return
	}

	self.drawHeaders()

	self.servePaintCellRangeRequest(self.startColumn, self.startRow, self.endColumn, self.endRow)

	self.servePaintSelectionRequest()
}

func (self *Sheet) servePaintCellRangeRequest(colStart int64, rowStart int64, colEnd int64, rowEnd int64) {

	if self == nil {
		return
	}

	c1, r1, c2, r2 := self.trimRangeToView(colStart, rowStart, colEnd, rowEnd)

	self.drawRange(c1, r1, c2, r2)
}

// Warning : assumes self.mark is well-ordered
func (self *Sheet) servePaintSelectionRequest() {

	if self == nil {
		return
	}

	// check if mark is out of view
	if self.mark.C1 > self.endColumn || self.mark.C2 < self.startColumn || self.mark.R1 > self.endRow || self.mark.R2 < self.startRow {
		return
	}

	c1, r1, c2, r2 := self.trimRangeToView(self.mark.C1, self.mark.R1, self.mark.C2, self.mark.R2)
	_, _, _, _, xlow, xhigh, ylow, yhigh := self.getIndicesAndRect(c1, r1, c2, r2)

	if !self.mark.IsSingleCell() {
		strokeFillRect(self.canvasContext, xlow, ylow, xhigh, yhigh, SELECTION_STROKE_COLOR, SELECTION_FILL_COLOR)
	}

	// !!!!!!!!!!!!!!!!!! TODO: For first top-left cell, draw thick border and last cell (bottom-right)
}

func (self *Sheet) drawHeaders() {

	if self == nil {
		return
	}
	// TODO: draw headers in the same way we did drawRange.
}

// Warning : no limit check for args here !
func (self *Sheet) drawRange(c1, r1, c2, r2 int64) {

	if self == nil {
		return
	}

	startXIdx, endXIdx, startYIdx, endYIdx, xlow, xhigh, ylow, yhigh := self.getIndicesAndRect(c1, r1, c2, r2)

	// cleanup the cell-range area
	noStrokeFillRect(self.canvasContext, xlow, ylow, xhigh, yhigh, CELL_DEFAULT_FILL_COLOR)

	// draw N vertical lines where N is number of columns in the range
	drawVertLines(self.canvasContext, self.colStartXCoords[startXIdx:endXIdx+1], ylow, yhigh, GRID_LINE_COLOR)

	// draw last vertical line to show end of last column
	drawLine(self.canvasContext, xhigh, ylow, xhigh, yhigh, GRID_LINE_COLOR)

	// draw N horizontal lines where N is number of rows in the range
	drawHorizLines(self.canvasContext, self.rowStartYCoords[startYIdx:endYIdx+1], xlow, xhigh, GRID_LINE_COLOR)

	// draw last horizontal line to show end of last row
	drawLine(self.canvasContext, xlow, yhigh, xhigh, yhigh, GRID_LINE_COLOR)

	self.drawCellRangeContents(c1, r1, c2, r2)

}

func (self *Sheet) drawCellRangeContents(c1, r1, c2, r2 int64) {

	startXIdx, endXIdx, startYIdx, endYIdx := self.getIndices(c1, r1, c2, r2)

	for cidx := startXIdx; cidx <= endXIdx; cidx++ {
		for ridx := startYIdx; ridx <= endYIdx; ridx++ {
			// TODO: write a single cell content drawing function. pass (cellstartX, cellstartY, cellWidth, cellHeight)
		}
	}
}

// Warning : no limit checks for args here !
func (self *Sheet) getIndices(c1, r1, c2, r2 int64) (startXIdx, endXIdx, startYIdx, endYIdx int64) {

	// index of start cell and end cell
	startXIdx = c1 - self.startColumn
	endXIdx = c2 - self.startColumn
	// index of start cell and end cell
	startYIdx = r1 - self.startRow
	endYIdx = r2 - self.startRow

	return
}

// Warning : no limit checks for args here !
func (self *Sheet) getIndicesAndRect(c1, r1, c2, r2 int64) (startXIdx, endXIdx, startYIdx, endYIdx int64,
	xlow, xhigh, ylow, yhigh float64) {

	startXIdx, endXIdx, startYIdx, endYIdx = getIndices(c1, r1, c2, r2)

	xlow = self.colStartXCoords[startXIdx]
	xhigh = math.Min(self.colStartXCoords[endXIdx+1], self.maxX) // end of last column in view

	ylow = self.rowStartYCoords[startYIdx]
	yhigh = math.Min(self.rowStartYCoords[endYIdx+1], self.maxY) // end of last row in view

	return
}

func (self *Sheet) trimRangeToView(colStart int64, rowStart int64, colEnd int64, rowEnd int64) (c1, r1, c2, r2 int64) {

	return maxInt64(colStart, self.startColumn), maxInt64(rowStart, self.startRow),
		minInt64(colEnd, self.endColumn), minInt64(rowEnd, self.endRow)
}

func (self *Sheet) addPaintRequest(request *SheetPaintRequest) {

	if self == nil {
		return
	}
	self.paintQueue <- request
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
