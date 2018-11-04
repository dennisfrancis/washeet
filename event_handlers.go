package washeet

import (
	//	"fmt"
	"math"
	"syscall/js"
)

func (self *Sheet) setupMouseHandlers() {
	if self == nil {
		return
	}

	self.setupMousedownHandler()
	self.setupMouseupHandler()
	self.setupMousemoveHandler()
}

func (self *Sheet) teardownMouseHandlers() {
	if self == nil {
		return
	}

	self.teardownMousemoveHandler()
	self.teardownMouseupHandler()
	self.teardownMousedownHandler()
}

func (self *Sheet) setupMousedownHandler() {
	if self == nil {
		return
	}

	self.mousedownHandler = js.NewCallback(func(args []js.Value) {
		event := args[0]
		x := event.Get("offsetX").Float()
		y := event.Get("offsetY").Float()
		//fmt.Printf("click at (%f, %f)\n", x, y)
		buttonCode := event.Get("button").Int()
		if buttonCode == 0 {
			xi, yi := self.getCellIndex(x, y)
			//fmt.Printf("cell index = (%d, %d)\n", xi, yi)
			if xi < 0 || yi < 0 {
				return
			}
			self.mouseState.setLeftDown()
			col, row := self.startColumn+xi, self.startRow+yi
			self.selectionState.setRefStartCell(col, row)
			self.selectionState.setRefCurrCell(col, row)
			self.PaintCellSelection(col, row)
		} else if buttonCode == 2 {
			self.mouseState.setRightDown()
		}
	})

	self.canvasElement.Call("addEventListener", "mousedown", self.mousedownHandler)
}

func (self *Sheet) setupMouseupHandler() {
	if self == nil {
		return
	}

	self.mouseupHandler = js.NewCallback(func(args []js.Value) {
		event := args[0]
		buttonCode := event.Get("button").Int()
		if buttonCode == 0 {
			self.mouseState.setLeftUp()
		} else if buttonCode == 2 {
			self.mouseState.setRightUp()
		}
		/*
			x := event.Get("offsetX").Float()
			y := event.Get("offsetY").Float()
			fmt.Printf("mouseup at (%f, %f)\n", x, y)
		*/
	})

	self.canvasElement.Call("addEventListener", "mouseup", self.mouseupHandler)
}

func (self *Sheet) setupMousemoveHandler() {
	if self == nil {
		return
	}

	self.mousemoveHandler = js.NewCallback(func(args []js.Value) {
		event := args[0]
		x := event.Get("offsetX").Float()
		y := event.Get("offsetY").Float()
		xidx, yidx := self.getCellIndex(x, y)
		bx, by, nearestxidx, nearestyidx := self.getNearestBorderXY(x, y, xidx, yidx)

		// bx and by are the nearest cell's start coordinates
		// so should not show resize mouse pointer for start borders of first column(col-resize) or first row(row-resize)
		if math.Abs(x-bx) <= 1.0 && nearestxidx >= 1 && nearestyidx == -1 {
			self.canvasElement.Get("style").Set("cursor", "col-resize")
		} else if math.Abs(y-by) <= 1.0 && nearestyidx >= 1 && nearestxidx == -1 {
			self.canvasElement.Get("style").Set("cursor", "row-resize")
		} else if x >= self.origX && x <= self.maxX && y >= self.origY && y <= self.maxY {
			self.canvasElement.Get("style").Set("cursor", "cell")

			// selection of a range while in a drag operation
			if self.mouseState.isLeftDown() {
				self.ehMutex.Lock()
				defer self.ehMutex.Unlock()
				refStartCell := self.selectionState.getRefStartCell()
				refCurrCell := self.selectionState.getRefCurrCell()
				col, row := self.startColumn+xidx, self.startRow+yidx

				if refCurrCell.Col == col && refCurrCell.Row == row {
					return
				}

				refCurrCell.Col, refCurrCell.Row = col, row
				c1, c2 := getInOrder(refStartCell.Col, col)
				r1, r2 := getInOrder(refStartCell.Row, row)
				self.PaintCellRangeSelection(c1, r1, c2, r2)
			}
		} else {
			// for headers
			self.canvasElement.Get("style").Set("cursor", "auto")
		}
	})

	self.canvasElement.Call("addEventListener", "mousemove", self.mousemoveHandler)
}

func (self *Sheet) teardownMousedownHandler() {

	self.canvasElement.Call("removeEventListener", "mousedown", self.mousedownHandler)
	self.mousedownHandler.Release()
}

func (self *Sheet) teardownMouseupHandler() {

	self.canvasElement.Call("removeEventListener", "mouseup", self.mouseupHandler)
	self.mouseupHandler.Release()
}

func (self *Sheet) teardownMousemoveHandler() {

	self.canvasElement.Call("removeEventListener", "mousemove", self.mousemoveHandler)
	self.canvasElement.Get("style").Set("cursor", "auto")
	self.mousemoveHandler.Release()
}

func (self *Sheet) setupKeyboardHandlers() {
	if self == nil {
		return
	}

	self.setupKeydownHandler()
}

func (self *Sheet) teardownKeyboardHandlers() {
	if self == nil {
		return
	}

	self.teardownKeydownHandler()
}

func (self *Sheet) setupKeydownHandler() {
	if self == nil {
		return
	}

	// NewEventCallback is used so that we can prevent window scrolling. event.preventDefault calls does not work
	self.keydownHandler = js.NewEventCallback(js.PreventDefault|js.StopPropagation|js.StopImmediatePropagation, func(event js.Value) {
		self.ehMutex.Lock()
		defer self.ehMutex.Unlock()
		keycode := event.Get("keyCode").Int()
		shiftKeyDown := event.Get("shiftKey").Bool()
		ctrlKeyDown := event.Get("ctrlKey").Bool()
		if keycode >= 37 && keycode <= 40 {
			self.arrowKeyHandler(keycode, shiftKeyDown)
		} else if ctrlKeyDown {
			// some command like Ctrl+C, Ctrl+V
			self.keyboardCommandHandler(keycode)
		}
	})

	self.canvasElement.Call("addEventListener", "keydown", self.keydownHandler)
}

func (self *Sheet) teardownKeydownHandler() {
	if self == nil {
		return
	}

	self.canvasElement.Call("removeEventListener", "keydown", self.keydownHandler)
	self.keydownHandler.Release()
}

func (self *Sheet) arrowKeyHandler(keycode int, shiftKeyDown bool) {
	if self == nil {
		return
	}

	refStartCell := self.selectionState.getRefStartCell()
	refCurrCell := self.selectionState.getRefCurrCell()
	col, row := refStartCell.Col, refStartCell.Row
	if shiftKeyDown {
		col, row = refCurrCell.Col, refCurrCell.Row
	}
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
	if paintFlag {
		if shiftKeyDown {
			self.selectionState.setRefCurrCell(col, row)
			c1, c2 := getInOrder(col, refStartCell.Col)
			r1, r2 := getInOrder(row, refStartCell.Row)
			self.PaintCellRangeSelection(c1, r1, c2, r2)
		} else {
			// Change both startCell and currCell references
			self.selectionState.setRefStartCell(col, row)
			self.selectionState.setRefCurrCell(col, row)
			self.PaintCellSelection(col, row)
		}
	}
}

func (self *Sheet) keyboardCommandHandler(keycode int) {
	// Ctrl key is down too, thats why this is now a command.
	if keycode == int('c') || keycode == int('C') {
		// Ctrl+C
		self.copySelectionToClipboard()
	}
}
