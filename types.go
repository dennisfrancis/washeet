package washeet

import (
	"sync"
	"syscall/js"
)

type sheetPaintType byte

type sheetPaintRequest struct {
	kind   sheetPaintType
	col    int64
	row    int64
	endCol int64
	endRow int64
	// For whole sheet paint requests
	changeSheetStartCol bool
	changeSheetStartRow bool
}

// TextAlignType is used to represent alignment of a cell's content.
type TextAlignType byte

// SheetDataProvider is the interface that needs to be implemented by the user of washeet
// which is used to draw and populate the contents of the spreadsheet.
type SheetDataProvider interface {
	// GetDisplayString returns the content of the cell at (column,row) as a string.
	GetDisplayString(column int64, row int64) string

	// GetColumnWidth returns the width of "column" column in pixels.
	GetColumnWidth(column int64) float64

	// GetRowHeight returns the height of "row" row in pixels.
	GetRowHeight(row int64) float64

	// TrimToNonEmptyRange trims given range represented by
	// { top-left cell (column = c1, row = r1), bottom-right cell (column = c2, row = r2)
	// to the biggest sub-range that does not have any leading/trailing empty columns/rows.
	// It returns false if given range is completely empty, else it returns true.
	TrimToNonEmptyRange(c1, r1, c2, r2 *int64) bool
}

// SheetModelUpdater is the interface that needs to be implemented by the user of washeet
// which is used to let communicate the changes in the contents of the spreadsheet..
type SheetModelUpdater interface {
	// SetColumnWidth updates the specified column's width in the spreadsheet model implementer.
	SetColumnWidth(column int64, width float64)

	// SetColumnWidth updates the specified rows's height in the spreadsheet model implementer.
	SetRowHeight(row int64, height float64)

	// SetCellContent updates the specified cell's content in the spreadsheet model implementer.
	SetCellContent(row, column int64, content string)
}

// markData represents a selection range as
// {top-left-cell(column = C1, row = R1), bottom-right-cell(column = C2, row = R2) }
type markData struct {
	c1 int64
	r1 int64
	c2 int64
	r2 int64
}

// cellCoords represents a cell's absolute coordinates.
type cellCoords struct {
	col int64
	row int64
}

// mouseState represents the state of the mouse at a given time.
type mouseState byte

// selectionState is used to store the start and current location of an on-going selection.
type selectionState struct {
	refStartCell cellCoords
	refCurrCell  cellCoords
}

type layoutData struct {
	startColumn int64
	startRow    int64

	endColumn int64
	endRow    int64

	colStartXCoords []float64
	rowStartYCoords []float64

	layoutFromStartCol bool
	layoutFromStartRow bool
}

// Sheet represents the spreadsheet user-interface.
type Sheet struct {
	document          js.Value
	window            js.Value
	navigator         js.Value
	container         *js.Value
	canvasElement     *js.Value
	canvasContext     js.Value
	clipboardTextArea js.Value
	origX             float64
	origY             float64
	maxX              float64
	maxY              float64

	dataSource SheetDataProvider
	dataSink   SheetModelUpdater

	rafLayoutData      *layoutData
	evtHndlrLayoutData *layoutData

	paintQueue        chan *sheetPaintRequest
	rafWorkerCallback js.Callback

	mark markData

	stopSignal  bool
	stopRequest chan struct{}

	ehMutex sync.Mutex

	mouseState     mouseState
	selectionState selectionState

	mousedownHandler js.Callback
	mouseupHandler   js.Callback
	mousemoveHandler js.Callback

	keydownHandler js.Callback

	layoutFromStartCol bool
	layoutFromStartRow bool
}

type cellMeasureGetter func(int64) float64

// Color represents a color.
type Color uint32

type colorSettings struct {
	gridLine        string
	cellFill        string
	cellStroke      string // fonts etc.
	cursorStroke    string
	headerFill      string
	selectionStroke string
	selectionFill   string
}
