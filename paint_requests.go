package washeet

import "fmt"

func (sheet *Sheet) addPaintRequest(request *sheetPaintRequest) bool {

	if sheet == nil || sheet.stopSignal {
		return false
	}

	queued := false
	select {
	case sheet.paintQueue <- request:
		queued = true
	default:
		// Queue is full, drop request
		fmt.Printf("[D]")
	}
	return queued
}

// if col/row = -1 no changes are made before whole-redraw
// changeSheetStartCol/changeSheetStartRow is also used to set sheet.layoutFromStartCol/sheet.layoutFromStartRow
func (sheet *Sheet) paintWholeSheet(col, row int64, changeSheetStartCol, changeSheetStartRow bool) bool {
	req := &sheetPaintRequest{
		kind:                sheetPaintWholeSheet,
		col:                 col,
		row:                 row,
		changeSheetStartCol: changeSheetStartCol,
		changeSheetStartRow: changeSheetStartRow,
	}

	if sheet.addPaintRequest(req) {
		// Layout will be partly/fully computed in RAF thread, so independently compute that here too.
		sheet.computeLayout(sheet.evtHndlrLayoutData, col, row, changeSheetStartCol, changeSheetStartRow)
		return true
	}

	return false
}

func (sheet *Sheet) paintCell(col int64, row int64) bool {

	if sheet == nil {
		return false
	}

	layout := sheet.evtHndlrLayoutData

	// optimization : don't fill the queue with these
	// if we know they are not going to be painted.
	if col < layout.startColumn || col > layout.endColumn ||
		row < layout.startRow || row > layout.endRow {
		return false
	}

	return sheet.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintCell,
		col:    col,
		row:    row,
		endCol: col,
		endRow: row,
	})
}

func (sheet *Sheet) paintCellRange(colStart int64, rowStart int64, colEnd int64, rowEnd int64) bool {

	if sheet == nil {
		return false
	}

	return sheet.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintCellRange,
		col:    colStart,
		row:    rowStart,
		endCol: colEnd,
		endRow: rowEnd,
	})
}

func (sheet *Sheet) paintCellSelection(col, row int64) bool {
	if sheet == nil {
		return false
	}

	return sheet.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintSelection,
		col:    col,
		row:    row,
		endCol: col,
		endRow: row,
	})
}

func (sheet *Sheet) paintCellRangeSelection(colStart, rowStart, colEnd, rowEnd int64) bool {
	if sheet == nil {
		return false
	}

	return sheet.addPaintRequest(&sheetPaintRequest{
		kind:   sheetPaintSelection,
		col:    colStart,
		row:    rowStart,
		endCol: colEnd,
		endRow: rowEnd,
	})
}
