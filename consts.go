package washeet

const (
	constMaxRowCount           int64   = 1048576
	constMaxRow                int64   = constMaxRowCount - 1
	constMaxColCount           int64   = 16384
	constMaxCol                int64   = constMaxColCount - 1
	constDefaultCellWidth      float64 = 90.0
	constDefaultCellHeight     float64 = 25.0
	constSheetPaintQueueLength int     = 50
)

const (
	sheetPaintWholeSheet sheetPaintType = iota
	sheetPaintCell
	sheetPaintCellRange
	sheetPaintSelection
)

const (
	// AlignLeft indicates left-alignment of cell-content.
	AlignLeft TextAlignType = iota
	// AlignCenter indicates center-alignment of cell-content.
	AlignCenter
	// AlignRight indicates right-alignment of cell-content.
	AlignRight
)
