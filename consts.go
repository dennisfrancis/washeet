package washeet

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

const (
	SheetPaintWholeSheet SheetPaintType = iota
	SheetPaintCell
	SheetPaintCellRange
	SheetPaintSelection
)

const (
	AlignLeft TextAlignType = iota
	AlignCenter
	AlignRight
)
