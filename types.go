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

type TextAlignType byte

type SheetDataProvider interface {
	GetDisplayString(column int64, row int64) string
	GetColumnWidth(column int64) float64
	GetRowHeight(row int64) float64
	// Trims given range to biggest range that does not have
	// leading/trailing empty columns/rows.
	// Returns false if given range is completely empty,
	// else returns true.
	TrimToNonEmptyRange(c1, r1, c2, r2 *int64) bool
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

type CellCoords struct {
	Col int64
	Row int64
}

type MouseState byte

type SelectionState struct {
	refStartCell CellCoords
	refCurrCell  CellCoords
}

type layoutData struct {
	startColumn int64
	startRow    int64

	endColumn int64
	endRow    int64

	colStartXCoords []float64
	rowStartYCoords []float64
}

type Sheet struct {
	document          js.Value
	window            js.Value
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

	startColumn int64
	startRow    int64

	endColumn int64
	endRow    int64

	paintQueue        chan *sheetPaintRequest
	rafWorkerCallback js.Callback

	colStartXCoords []float64
	rowStartYCoords []float64

	mark MarkData

	stopSignal  bool
	stopRequest chan struct{}

	ehMutex sync.Mutex

	mouseState     MouseState
	selectionState SelectionState

	mousedownHandler js.Callback
	mouseupHandler   js.Callback
	mousemoveHandler js.Callback

	keydownHandler js.Callback

	layoutFromStartCol bool
	layoutFromStartRow bool
}

type cellMeasureGetter func(int64) float64
