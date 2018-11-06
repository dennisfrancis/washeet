package washeet

const (
	MAXROWCOUNT              int64   = 1048576
	MAXROW                   int64   = MAXROWCOUNT - 1
	MAXCOLCOUNT              int64   = 16384
	MAXCOL                   int64   = MAXCOLCOUNT - 1
	DEFAULT_CELL_WIDTH       float64 = 90.0
	DEFAULT_CELL_HEIGHT      float64 = 25.0
	SHEET_PAINT_QUEUE_LENGTH int     = 50

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
