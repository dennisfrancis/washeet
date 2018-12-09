package washeet

func (sheet *Sheet) setupCanvas(width, height float64) {

	sheet.container.Call("setAttribute", "style", "position:relative;margin:0;padding:0;")
	sheet.canvasElement = sheet.document.Call("createElement", "canvas")
	sheet.canvasElement.Set("width", width)
	sheet.canvasElement.Set("height", height)
	sheet.container.Call("appendChild", sheet.canvasElement)
	// Set style to control positioning and z-index.
	sheet.canvasElement.Call("setAttribute", "style", "position:absolute;z-index:1;margin:0;padding:0;outline:none;")
	// Make it the default focussed element.
	sheet.canvasElement.Set("tabIndex", 1)
	sheet.canvasContext = sheet.canvasElement.Call("getContext", "2d")
}

func (sheet *Sheet) teardownCanvas() {

	sheet.container.Call("removeChild", sheet.canvasElement)
}
