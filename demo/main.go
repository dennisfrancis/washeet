package main

import (
	"fmt"
	"math"
	"syscall/js"

	"github.com/dennisfrancis/washeet"
)

type SheetModel struct {
}

var (
	redColor   *washeet.Color = washeet.NewColor(200, 0, 0)
	greenColor *washeet.Color = washeet.NewColor(0, 200, 0)
	blueColor  *washeet.Color = washeet.NewColor(0, 0, 200)

	lredColor   *washeet.Color = washeet.NewColor(255, 242, 230)
	lgreenColor *washeet.Color = washeet.NewColor(204, 255, 204)
	lblueColor  *washeet.Color = washeet.NewColor(230, 255, 255)
)

// Satisfy SheetDataProvider interface.

func (self *SheetModel) GetDisplayString(column int64, row int64) string {
	if column == 5 && row == 5 {
		return "Hello washeet !"
	}

	return fmt.Sprintf("Cell(%d, %d)", column, row)
}

func (self *SheetModel) GetCellAttribs(column, row int64) *washeet.CellAttribs {
	idx := column % int64(4)
	attrib := washeet.NewDefaultCellAttribs()
	if idx == 1 {
		attrib.SetBold(true)
	} else if idx == 2 {
		attrib.SetItalics(true)
	} else if idx == 3 {
		attrib.SetUnderline(true)
	}

	idxrow := row % int64(9)
	if idxrow < 3 {
		attrib.SetAlignment(washeet.AlignLeft)
		attrib.SetFGColor(redColor)
	} else if idxrow < 6 {
		attrib.SetAlignment(washeet.AlignCenter)
		attrib.SetFGColor(greenColor)
		attrib.SetBGColor(lgreenColor)
	} else {
		attrib.SetAlignment(washeet.AlignRight)
		attrib.SetFGColor(blueColor)
		attrib.SetBGColor(lblueColor)
	}

	return attrib
}

func (self *SheetModel) GetColumnWidth(column int64) float64 {
	if column == 5 {
		return 230.0
	}
	return 0.0
}

func (self *SheetModel) GetRowHeight(row int64) float64 {
	if row == 5 {
		return 70.0
	}
	return 0.0
}

func (self *SheetModel) TrimToNonEmptyRange(c1, r1, c2, r2 *int64) bool {
	return true
}

func main() {

	fmt.Println("Hello washeet !")

	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "washeet")
	container := doc.Call("getElementById", "container")
	width, height := doc.Get("body").Get("clientWidth").Float(), doc.Get("body").Get("clientHeight").Float()
	width, height = math.Floor(width*0.85), math.Floor(height*0.85)
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	closeButton := doc.Call("getElementById", "close-button")
	quit := make(chan bool)

	closeHandler := js.NewCallback(func(args []js.Value) {
		quit <- true
	})
	closeButton.Call("addEventListener", "click", closeHandler)

	model := &SheetModel{}
	sheet := washeet.NewSheet(&canvasEl, &container, 0.0, 0.0, width-1.0, height-1.0, model)
	sheet.Start()

	<-quit
	closeButton.Call("removeEventListener", "click", closeHandler)
	closeHandler.Release()
	sheet.Stop()
	fmt.Println("sheet.Stop() successfull !")
}
