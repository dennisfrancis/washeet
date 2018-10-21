package washeet

import (
	"syscall/js"
)

type SheetPaintType byte

type SheetPaintRequest struct {
	Kind   SheetPaintType
	Col    int64
	Row    int64
	EndCol int64
	EndRow int64
}

type TextAlignType byte

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

type MouseState byte

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

	mouseState       MouseState
	mousedownHandler js.Callback
	mouseupHandler   js.Callback
	mousemoveHandler js.Callback
}
