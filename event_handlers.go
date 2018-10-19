package washeet

import (
	"math"
	"syscall/js"
)

func (self *Sheet) setupClickHandler() {
	if self == nil {
		return
	}

	self.clickHandler = js.NewCallback(func(args []js.Value) {
		event := args[0]
		x := event.Get("offsetX").Float()
		y := event.Get("offsetY").Float()
		//fmt.Printf("click at (%f, %f)\n", x, y)

		xi, yi := self.getCellIndex(x, y)
		//fmt.Printf("cell index = (%d, %d)\n", xi, yi)
		if xi < 0 || yi < 0 {
			return
		}
		currsel := &(self.mark)
		self.PaintCellRange(currsel.C1, currsel.R1, currsel.C2, currsel.R2)
		self.PaintCellSelection(self.startColumn+xi, self.startRow+yi)
	})

	self.canvasElement.Call("addEventListener", "click", self.clickHandler)
}

func (self *Sheet) setupMousemoveHandler() {
	if self == nil {
		return
	}

	self.mousemoveHandler = js.NewCallback(func(args []js.Value) {
		event := args[0]
		x := event.Get("offsetX").Float()
		y := event.Get("offsetY").Float()
		bx, by, cellxidx, cellyidx := self.getNearestBorderXY(x, y)

		// bx and by are the nearest cell's start coordinates
		// so should not show resize mouse pointer for start borders of first column(col-resize) or first row(row-resize)
		if math.Abs(x-bx) <= 1.0 && cellxidx >= 1 && cellyidx == -1 {
			self.canvasElement.Get("style").Set("cursor", "col-resize")
		} else if math.Abs(y-by) <= 1.0 && cellyidx >= 1 && cellxidx == -1 {
			self.canvasElement.Get("style").Set("cursor", "row-resize")
		} else if x >= self.origX && x <= self.maxX && y >= self.origY && y <= self.maxY {
			self.canvasElement.Get("style").Set("cursor", "cell")
		} else {
			// for headers
			self.canvasElement.Get("style").Set("cursor", "auto")
		}
	})

	self.canvasElement.Call("addEventListener", "mousemove", self.mousemoveHandler)
}
