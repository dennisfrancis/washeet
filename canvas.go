package washeet

func (sheet *Sheet) setupCanvas(width, height float64) {

	sheet.container.Call("setAttribute", "style", "position:relative;margin:0;padding:0;")
	sheet.canvasStore = canvasStoreType{
		sheetCanvasElement:     sheet.document.Call("createElement", "canvas"),
		selectionCanvasElement: sheet.document.Call("createElement", "canvas"),
	}

	// TODO: avoid following code duplication
	// Background canvas :- the sheet canvas
	sheetCanvasElement := sheet.canvasStore.sheetCanvasElement
	sheetCanvasElement.Set("width", width)
	sheetCanvasElement.Set("height", height)
	sheet.container.Call("appendChild", sheetCanvasElement)
	// Set style to control positioning and z-index.
	sheetCanvasElement.Call("setAttribute", "style", "position:absolute;left:0;top:0;z-index:1;margin:0;padding:0;outline:none;")

	selectionCanvasElement := sheet.canvasStore.selectionCanvasElement
	selectionCanvasElement.Set("width", width)
	selectionCanvasElement.Set("height", height)
	sheet.container.Call("appendChild", selectionCanvasElement)
	// Set style to control positioning and z-index.
	selectionCanvasElement.Call("setAttribute", "style", "position:absolute;left:0;top:0;z-index:2;margin:0;padding:0;outline:none;")

	// Make selectionCanvasElement the default focussed element.
	selectionCanvasElement.Set("tabIndex", 1)
	sheet.canvasStore.sheetCanvasContext = sheetCanvasElement.Call("getContext", "2d")
	sheet.canvasStore.selectionCanvasContext = selectionCanvasElement.Call("getContext", "2d")

	sheet.canvasStore.foregroundCanvasElement = selectionCanvasElement
}

func (sheet *Sheet) teardownCanvas() {

	sheet.container.Call("removeChild", sheet.canvasStore.selectionCanvasElement)
	sheet.container.Call("removeChild", sheet.canvasStore.sheetCanvasElement)
}
