package main

import (
	"fmt"
	"github.com/dennisfrancis/washeet"
	"syscall/js"
)

var (
	width  float64
	height float64
	ctx    js.Value
)

type SheetModel struct {
}

// Satisfy SheetDataProvider interface.

func (self *SheetModel) GetDisplayString(column int64, row int64) string {
	if column == 5 && row == 5 {
		return "Hello washeet !"
	}

	return ""
}

func (self *SheetModel) GetColumnWidth(column int64) float64 {
	if column == 5 {
		return 100.0
	}
	return 0.0
}

func (self *SheetModel) GetRowHeight(row int64) float64 {
	return 0.0
}

// Satisfy SheetModelUpdater interface.

func (self *SheetModel) SetColumnWidth(column int64, width float64) {
}

func (self *SheetModel) SetRowHeight(row int64, height float64) {
}

func (self *SheetModel) SetCellContent(row, column int64, content string) {
}

func main() {

	fmt.Println("Hello washeet !")

	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "washeet")
	width = doc.Get("body").Get("clientWidth").Float()
	height = doc.Get("body").Get("clientHeight").Float()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	ctx = canvasEl.Call("getContext", "2d")
	closeButton := doc.Call("getElementById", "close-button")
	quit := make(chan bool)

	closeHandler := js.NewCallback(func(args []js.Value) {
		quit <- true
	})
	closeButton.Call("addEventListener", "click", closeHandler)

	model := &SheetModel{}
	sheet := washeet.NewSheet(&ctx, 0.0, 0.0, width-1.0, height-1.0, model, model)
	sheet.Start()

	<-quit
	closeButton.Call("removeEventListener", "click", closeHandler)
	closeHandler.Release()
	sheet.Stop()
	fmt.Println("sheet.Stop() successfull !")
}
