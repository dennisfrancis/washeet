package washeet

import (
	//"fmt"
	"math"
	"syscall/js"
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

	sheet.mousedownHandler = js.NewCallback(func(args []js.Value) {
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
				return
			}
			sheet.mouseState.setLeftDown()
			col, row := layout.startColumn+xi, layout.startRow+yi
			sheet.selectionState.setRefStartCell(col, row)
			sheet.selectionState.setRefCurrCell(col, row)
			sheet.PaintCellSelection(col, row)
		} else if buttonCode == 2 {
			sheet.mouseState.setRightDown()
		}
	})

	sheet.canvasElement.Call("addEventListener", "mousedown", sheet.mousedownHandler)
}

func (sheet *Sheet) setupMouseupHandler() {
	if sheet == nil {
		return
	}

	sheet.mouseupHandler = js.NewCallback(func(args []js.Value) {
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
	})

	sheet.canvasElement.Call("addEventListener", "mouseup", sheet.mouseupHandler)
}

func (sheet *Sheet) setupMousemoveHandler() {
	if sheet == nil {
		return
	}

	sheet.mousemoveHandler = js.NewCallback(func(args []js.Value) {
		event := args[0]
		x := event.Get("offsetX").Float()
		y := event.Get("offsetY").Float()
		layout := sheet.evtHndlrLayoutData
		xidx, yidx := sheet.getCellIndex(layout, x, y)
		bx, by, nearestxidx, nearestyidx := sheet.getNearestBorderXY(layout, x, y, xidx, yidx)

		// bx and by are the nearest cell's start coordinates
		// so should not show resize mouse pointer for start borders of first column(col-resize) or first row(row-resize)
		if math.Abs(x-bx) <= 1.0 && nearestxidx >= 1 && nearestyidx == -1 {
			sheet.canvasElement.Get("style").Set("cursor", "col-resize")
		} else if math.Abs(y-by) <= 1.0 && nearestyidx >= 1 && nearestxidx == -1 {
			sheet.canvasElement.Get("style").Set("cursor", "row-resize")
		} else if x >= sheet.origX && x <= sheet.maxX && y >= sheet.origY && y <= sheet.maxY {
			sheet.canvasElement.Get("style").Set("cursor", "cell")

			// selection of a range while in a drag operation
			if sheet.mouseState.isLeftDown() {
				sheet.ehMutex.Lock()
				defer sheet.ehMutex.Unlock()
				refStartCell := sheet.selectionState.getRefStartCell()
				refCurrCell := sheet.selectionState.getRefCurrCell()
				col, row := layout.startColumn+xidx, layout.startRow+yidx

				if refCurrCell.Col == col && refCurrCell.Row == row {
					return
				}

				refCurrCell.Col, refCurrCell.Row = col, row
				c1, c2 := getInOrder(refStartCell.Col, col)
				r1, r2 := getInOrder(refStartCell.Row, row)
				sheet.PaintCellRangeSelection(c1, r1, c2, r2)
			}
		} else {
			// for headers
			sheet.canvasElement.Get("style").Set("cursor", "auto")
		}
	})

	sheet.canvasElement.Call("addEventListener", "mousemove", sheet.mousemoveHandler)
}

func (sheet *Sheet) teardownMousedownHandler() {

	sheet.canvasElement.Call("removeEventListener", "mousedown", sheet.mousedownHandler)
	sheet.mousedownHandler.Release()
}

func (sheet *Sheet) teardownMouseupHandler() {

	sheet.canvasElement.Call("removeEventListener", "mouseup", sheet.mouseupHandler)
	sheet.mouseupHandler.Release()
}

func (sheet *Sheet) teardownMousemoveHandler() {

	sheet.canvasElement.Call("removeEventListener", "mousemove", sheet.mousemoveHandler)
	sheet.canvasElement.Get("style").Set("cursor", "auto")
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
	sheet.keydownHandler = js.NewEventCallback(js.PreventDefault|js.StopPropagation|js.StopImmediatePropagation, func(event js.Value) {
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
	})

	sheet.canvasElement.Call("addEventListener", "keydown", sheet.keydownHandler)
}

func (sheet *Sheet) teardownKeydownHandler() {
	if sheet == nil {
		return
	}

	sheet.canvasElement.Call("removeEventListener", "keydown", sheet.keydownHandler)
	sheet.keydownHandler.Release()
}

func (sheet *Sheet) arrowKeyHandler(keycode int, shiftKeyDown bool) {
	if sheet == nil {
		return
	}

	refStartCell := sheet.selectionState.getRefStartCell()
	refCurrCell := sheet.selectionState.getRefCurrCell()
	col, row := refStartCell.Col, refStartCell.Row
	if shiftKeyDown {
		col, row = refCurrCell.Col, refCurrCell.Row
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
		if col == MAXCOL {
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
		if row == MAXROW {
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
		sheet.PaintWholeSheet(changeCol, changeRow, changeSheetStartCol, changeSheetStartRow)
	}

	if paintFlag {
		if shiftKeyDown {
			sheet.selectionState.setRefCurrCell(col, row)
			c1, c2 := getInOrder(col, refStartCell.Col)
			r1, r2 := getInOrder(row, refStartCell.Row)
			sheet.PaintCellRangeSelection(c1, r1, c2, r2)
		} else {
			// Change both startCell and currCell references
			sheet.selectionState.setRefStartCell(col, row)
			sheet.selectionState.setRefCurrCell(col, row)
			sheet.PaintCellSelection(col, row)
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
