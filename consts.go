package washeet

const (
	constMaxRowCount           int64   = 1048576
	constMaxRow                int64   = constMaxRowCount - 1
	constMaxColCount           int64   = 16384
	constMaxCol                int64   = constMaxColCount - 1
	constDefaultCellWidth      float64 = 90.0
	constDefaultCellHeight     float64 = 25.0
	constSheetPaintQueueLength int     = 50

	GRID_LINE_COLOR           string = "rgba(200, 200, 200, 1.0)"
	CELL_DEFAULT_FILL_COLOR   string = "rgba(255, 255, 255, 1.0)"
	CELL_DEFAULT_STROKE_COLOR string = "rgba(0, 0, 0, 1.0)" // fonts etc.
	CURSOR_STROKE_COLOR       string = "rgba(0, 0, 0, 1.0)"
	HEADER_FILL_COLOR         string = "rgba(240, 240, 240, 1.0)"
	SELECTION_STROKE_COLOR    string = "rgba(100, 156, 244, 1.0)"
	SELECTION_FILL_COLOR      string = "rgba(100, 156, 244, 0.2)"
)

const (
	sheetPaintWholeSheet sheetPaintType = iota
	sheetPaintCell
	sheetPaintCellRange
	sheetPaintSelection
)

const (
	AlignLeft TextAlignType = iota
	AlignCenter
	AlignRight
)
