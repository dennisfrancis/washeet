package washeet

import (
	"syscall/js"
)

var path2dCtor = js.Global().Get("Path2D")

func drawVertLines(canvasContext *js.Value, xcoords []float64, ylow, yhigh float64, color string) {
	path := path2dCtor.New()
	canvasContext.Set("strokeStyle", color)
	for _, xcoord := range xcoords {
		addLine2Path(&path, xcoord, ylow, xcoord, yhigh)
	}
	canvasContext.Call("stroke", path)
}

func drawHorizLines(canvasContext *js.Value, ycoords []float64, xlow, xhigh float64, color string) {
	path := path2dCtor.New()
	canvasContext.Set("strokeStyle", color)
	for _, ycoord := range ycoords {
		addLine2Path(&path, xlow, ycoord, xhigh, ycoord)
	}
	canvasContext.Call("stroke", path)
}

func drawLine(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, color string) {
	path := path2dCtor.New()
	canvasContext.Set("strokeStyle", color)
	addLine2Path(&path, xlow, ylow, xhigh, yhigh)
	canvasContext.Call("stroke", path)
}

// Adds a line to the path, no actual drawing takes place
func addLine2Path(path *js.Value, xlow, ylow, xhigh, yhigh float64) {
	path.Call("moveTo", xlow, ylow)
	path.Call("lineTo", xhigh, yhigh)
}

func noStrokeFillRect(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, fillColor string) {

}

func strokeFillRect(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, strokeColor, fillColor string) {

}

func strokeNoFillRect(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, strokeColor string) {

}

func drawText(canvasContext *js.Value, xlow, ylow, xhigh, yhigh, xmax, ymax float64, text string, align TextAlignType) {

}
