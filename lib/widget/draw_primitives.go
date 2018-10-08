package widget

import (
	"syscall/js"
)

func drawVertLines(canvasContext *js.Value, xcoords []float64, ylow, yhigh float64, color string) {
	for _, xcoord := range xcoords {
		drawLine(canvasContext, xcoord, ylow, xcoord, yhigh, color)
	}
}

func drawHorizLines(canvasContext *js.Value, ycoords []float64, xlow, xhigh float64, color string) {
	for _, ycoord := range ycoords {
		drawLine(canvasContext, xlow, ycoord, xhigh, ycoord, color)
	}
}

func drawLine(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, color string) {

}

func noStrokeFillRect(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, fillColor string) {

}

func strokeFillRect(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, strokeColor, fillColor string) {

}
