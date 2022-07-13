package washeet

import (
	"syscall/js"
	"time"
)

var (
	mouseMoveLast time.Time
)

func (sheet *Sheet) setupMouseHandlers() {
	if sheet == nil {
		return
	}

	sheet.setupMousedownHandler()
	sheet.setupMouseupHandler()
	sheet.setupMousemoveHandler()
}

func (sheet *Sheet) teardownMouseHandlers() {
	if sheet == nil {
		return
	}

	sheet.teardownMousemoveHandler()
	sheet.teardownMouseupHandler()
	sheet.teardownMousedownHandler()
}

func (sheet *Sheet) setupMousedownHandler() {
	if sheet == nil {
		return
	}

	sheet.mousedownHandler = js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		x := event.Get("offsetX").Float()
		y := event.Get("offsetY").Float()
		layout := sheet.evtHndlrLayoutData
		//fmt.Printf("click at (%f, %f)\n", x, y)
		buttonCode := event.Get("button").Int()
		if buttonCode == 0 {
			xi, yi := sheet.getCellIndex(layout, x, y)
			//fmt.Printf("cell index = (%d, %d)\n", xi, yi)
			if xi < 0 || yi < 0 {
				return nil
			}
			sheet.mouseState.setLeftDown()
			col, row := layout.startColumn+xi, layout.startRow+yi
			if sheet.paintCellSelection(col, row) {
				sheet.selectionState.setRefStartCell(col, row)
				sheet.selectionState.setRefCurrCell(col, row)
			}
		} else if buttonCode == 2 {
			sheet.mouseState.setRightDown()
		}
		return nil
	})

	sheet.canvasStore.foregroundCanvasElement.Call("addEventListener", "mousedown", sheet.mousedownHandler)
}

func (sheet *Sheet) setupMouseupHandler() {
	if sheet == nil {
		return
	}

	sheet.mouseupHandler = js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		buttonCode := event.Get("button").Int()
		if buttonCode == 0 {
			sheet.mouseState.setLeftUp()
		} else if buttonCode == 2 {
			sheet.mouseState.setRightUp()
		}
		/*
			x := event.Get("offsetX").Float()
			y := event.Get("offsetY").Float()
			fmt.Printf("mouseup at (%f, %f)\n", x, y)
		*/
		return nil
	})

	sheet.canvasStore.foregroundCanvasElement.Call("addEventListener", "mouseup", sheet.mouseupHandler)
	sheet.window.Call("addEventListener", "mouseup", sheet.mouseupHandler)
}

func (sheet *Sheet) setupMousemoveHandler() {
	if sheet == nil {
		return
	}

	sheet.mousemoveHandler = js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		x := event.Get("offsetX").Float()
		y := event.Get("offsetY").Float()
		layout := sheet.evtHndlrLayoutData
		xidx, yidx := sheet.getCellIndex(layout, x, y)
		cellgridOrigX, cellGridOrigY := sheet.origX+constDefaultCellWidth, sheet.origY+constDefaultCellHeight

		// bx, by, nearestxidx, nearestyidx := sheet.getNearestBorderXY(layout, x, y, xidx, yidx)
		// bx and by are the nearest cell's start coordinates
		// so should not show resize mouse pointer for start borders of first column(col-resize) or first row(row-resize)

		// We don't allow column/row resizing yet, so no point in showing these cursors.
		// if math.Abs(x-bx) <= 1.0 && nearestxidx >= 1 && nearestyidx == -1 {
		// 	sheet.canvasStore.foregroundCanvasElement.Get("style").Set("cursor", "col-resize")
		// } else if math.Abs(y-by) <= 1.0 && nearestyidx >= 1 && nearestxidx == -1 {
		// 	sheet.canvasStore.foregroundCanvasElement.Get("style").Set("cursor", "row-resize")
		// }

		if x >= cellgridOrigX && x <= sheet.maxX && y >= cellGridOrigY && y <= sheet.maxY {
			sheet.canvasStore.foregroundCanvasElement.Get("style").Set("cursor", "cell")

			// selection of a range while in a drag operation
			if sheet.mouseState.isLeftDown() {
				sheet.ehMutex.Lock()
				defer sheet.ehMutex.Unlock()
				refStartCell := sheet.selectionState.getRefStartCell()
				refCurrCell := sheet.selectionState.getRefCurrCell()
				col, row := layout.startColumn+xidx, layout.startRow+yidx

				if refCurrCell.col == col && refCurrCell.row == row {
					return nil
				}

				c1, c2 := getInOrder(refStartCell.col, col)
				r1, r2 := getInOrder(refStartCell.row, row)
				if sheet.paintCellRangeSelection(c1, r1, c2, r2) {
					refCurrCell.col, refCurrCell.row = col, row
				}

			}
		} else {
			// for headers
			sheet.canvasStore.foregroundCanvasElement.Get("style").Set("cursor", "auto")
			if time.Now().Sub(mouseMoveLast) < 100*time.Millisecond {
				return nil
			}
			mouseMoveLast = time.Now()
			if sheet.mouseState.isLeftDown() {
				sheet.ehMutex.Lock()
				defer sheet.ehMutex.Unlock()
				//fmt.Printf("xidx = %d ", xidx)
				changeCol, changeRow := int64(-1), int64(-1)
				changeColFromStart, changeRowFromStart := true, true
				refStartCell := sheet.selectionState.getRefStartCell()
				col, row := int64(-1), int64(-1)

				// Row change
				if y > sheet.maxY {
					changeRow = layout.endRow + 1
					changeRowFromStart = false
					row = changeRow
				} else if y < cellGridOrigY {
					if layout.startRow > 0 {
						changeRow = layout.startRow - 1
						changeRowFromStart = true
						row = changeRow
					}
				} else if yidx >= 0 {
					row = layout.startRow + yidx
				}

				// Column change
				if x > sheet.maxX {
					changeCol = layout.endColumn + 1
					changeColFromStart = false
					col = changeCol
				} else if x < cellgridOrigX {
					if layout.startColumn > 0 {
						changeCol = layout.startColumn - 1
						changeColFromStart = true
						col = changeCol
					}
				} else if xidx >= 0 {
					col = layout.startColumn + xidx
				}

				if col != int64(-1) || row != int64(-1) {
					c1, c2 := getInOrder(col, refStartCell.col)
					r1, r2 := getInOrder(row, refStartCell.row)
					if changeCol != int64(-1) || changeRow != int64(-1) {
						if sheet.paintWholeSheet(changeCol, changeRow, changeColFromStart, changeRowFromStart) {
							if sheet.paintCellRangeSelection(c1, r1, c2, r2) {
								sheet.selectionState.setRefCurrCell(col, row)
							}
						}
					} else if sheet.paintCellRangeSelection(c1, r1, c2, r2) {
						sheet.selectionState.setRefCurrCell(col, row)
					}
				}

			} // End of Mouse move + down - out of sheet cases handler
		}

		return nil
	})

	sheet.canvasStore.foregroundCanvasElement.Call("addEventListener", "mousemove", sheet.mousemoveHandler)
	sheet.window.Call("addEventListener", "mousemove", sheet.mousemoveHandler)
}

func (sheet *Sheet) teardownMousedownHandler() {

	sheet.canvasStore.foregroundCanvasElement.Call("removeEventListener", "mousedown", sheet.mousedownHandler)
	sheet.mousedownHandler.Release()
}

func (sheet *Sheet) teardownMouseupHandler() {

	sheet.canvasStore.foregroundCanvasElement.Call("removeEventListener", "mouseup", sheet.mouseupHandler)
	sheet.window.Call("removeEventListener", "mouseup", sheet.mouseupHandler)
	sheet.mouseupHandler.Release()
}

func (sheet *Sheet) teardownMousemoveHandler() {

	sheet.canvasStore.foregroundCanvasElement.Call("removeEventListener", "mousemove", sheet.mousemoveHandler)
	sheet.window.Call("removeEventListener", "mousemove", sheet.mousemoveHandler)
	sheet.canvasStore.foregroundCanvasElement.Get("style").Set("cursor", "auto")
	sheet.mousemoveHandler.Release()
}

func (sheet *Sheet) setupKeyboardHandlers() {
	if sheet == nil {
		return
	}

	sheet.setupKeydownHandler()
}

func (sheet *Sheet) teardownKeyboardHandlers() {
	if sheet == nil {
		return
	}

	sheet.teardownKeydownHandler()
}

func (sheet *Sheet) setupKeydownHandler() {
	if sheet == nil {
		return
	}

	// NewEventCallback is used so that we can prevent window scrolling. event.preventDefault calls does not work
	sheet.keydownHandler = js.FuncOf(func(this js.Value, args []js.Value) any {
		event := args[0]
		event.Call("preventDefault")
		event.Call("stopPropagation")
		event.Call("stopImmediatePropagation")
		sheet.ehMutex.Lock()
		defer sheet.ehMutex.Unlock()
		keycode := event.Get("keyCode").Int()
		shiftKeyDown := event.Get("shiftKey").Bool()
		ctrlKeyDown := event.Get("ctrlKey").Bool()
		if keycode >= 37 && keycode <= 40 {
			sheet.arrowKeyHandler(keycode, shiftKeyDown)
		} else if ctrlKeyDown {
			// some command like Ctrl+C, Ctrl+V
			sheet.keyboardCommandHandler(keycode)
		}
		return nil
	})

	sheet.canvasStore.foregroundCanvasElement.Call("addEventListener", "keydown", sheet.keydownHandler)
}

func (sheet *Sheet) teardownKeydownHandler() {
	if sheet == nil {
		return
	}

	sheet.canvasStore.foregroundCanvasElement.Call("removeEventListener", "keydown", sheet.keydownHandler)
	sheet.keydownHandler.Release()
}

func (sheet *Sheet) arrowKeyHandler(keycode int, shiftKeyDown bool) {
	if sheet == nil {
		return
	}

	refStartCell := sheet.selectionState.getRefStartCell()
	refCurrCell := sheet.selectionState.getRefCurrCell()
	col, row := refStartCell.col, refStartCell.row
	if shiftKeyDown {
		col, row = refCurrCell.col, refCurrCell.row
	}
	// flag for selection paint
	paintFlag := true
	switch keycode {
	case 37: // Left
		if col == 0 {
			paintFlag = false
		} else {
			col--
		}
	case 39: // Right
		if col == constMaxCol {
			paintFlag = false
		} else {
			col++
		}
	case 38: // Up
		if row == 0 {
			paintFlag = false
		} else {
			row--
		}
	case 40: // Down
		if row == constMaxRow {
			paintFlag = false
		} else {
			row++
		}

	}

	shouldPaintWhole := false
	changeCol, changeRow := int64(-1), int64(-1)
	changeSheetStartCol, changeSheetStartRow := true, true

	layout := sheet.evtHndlrLayoutData

	if col < layout.startColumn {
		changeCol = col
		shouldPaintWhole = true
	} else if col >= layout.endColumn {
		changeCol = col
		changeSheetStartCol = false
		shouldPaintWhole = true
	}

	if row < layout.startRow {
		changeRow = row
		shouldPaintWhole = true
	} else if row >= layout.endRow {
		changeRow = row
		changeSheetStartRow = false
		shouldPaintWhole = true
	}

	if shouldPaintWhole {
		sheet.paintWholeSheet(changeCol, changeRow, changeSheetStartCol, changeSheetStartRow)
	}

	if paintFlag {
		if shiftKeyDown {
			c1, c2 := getInOrder(col, refStartCell.col)
			r1, r2 := getInOrder(row, refStartCell.row)
			if sheet.paintCellRangeSelection(c1, r1, c2, r2) {
				sheet.selectionState.setRefCurrCell(col, row)
			}
		} else {
			if sheet.paintCellSelection(col, row) {
				// Change both startCell and currCell references
				sheet.selectionState.setRefStartCell(col, row)
				sheet.selectionState.setRefCurrCell(col, row)
			}
		}
	}
}

func (sheet *Sheet) keyboardCommandHandler(keycode int) {
	// Ctrl key is down too, thats why this is now a command.
	if keycode == int('c') || keycode == int('C') {
		// Ctrl+C
		//fmt.Println("Copy")
		sheet.copySelectionToClipboard()
	}
}
