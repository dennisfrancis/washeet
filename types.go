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

	// GetCellAttribs returns cell-attributes information via CellAttribs type.
	GetCellAttribs(column, row int64) *CellAttribs

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

// SheetNotifier is the interface via which client-users of washeet can
// notify which cell/range has changed its contents.
type SheetNotifier interface {
	// UpdateCell call tells washeet that the specified cell contents have changed
	// and needs redrawing.
	UpdateCell(column, row int64)

	// UpdateRange call tells washeet that the specified range's contents have changed
	// and needs redrawing.
	UpdateRange(startColumn, startRow, endColumn, endRow int64)
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

type canvasStoreType struct {
	sheetCanvasElement     js.Value
	sheetCanvasContext     js.Value
	selectionCanvasElement js.Value
	selectionCanvasContext js.Value

	foregroundCanvasElement js.Value
}

// Sheet represents the spreadsheet user-interface.
type Sheet struct {
	document          js.Value
	window            js.Value
	navigator         js.Value
	container         *js.Value
	canvasStore       canvasStoreType
	clipboardTextArea js.Value
	origX             float64
	origY             float64
	maxX              float64
	maxY              float64

	dataSource SheetDataProvider

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
	gridLine           *Color
	cellFill           *Color
	cellStroke         *Color // fonts etc.
	cursorStroke       *Color
	headerFill         *Color
	selectionStroke    *Color
	selectionFill      *Color
	selectionClearFill *Color
}

// textAttribs represent a cell's text attribute states like
// bold, italics, underline.
type textAttribs uint8

// CellAttribs stores the attributes of a cell.
type CellAttribs struct {
	txtAttribs *textAttribs
	txtAlign   TextAlignType
	fgColor    *Color
	bgColor    *Color
	fontSize   uint8
}

type cornerType uint8
