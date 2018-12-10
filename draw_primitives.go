package washeet

import (
	"fmt"
	"math"
	"syscall/js"
)

var path2dCtor = js.Global().Get("Path2D")

func drawVertLines(canvasContext *js.Value, xcoords []float64, ylow, yhigh float64, color *Color) {
	path := path2dCtor.New()
	canvasContext.Set("strokeStyle", color.toString())
	for _, xcoord := range xcoords {
		addLine2Path(&path, xcoord+0.5, ylow, xcoord+0.5, yhigh)
	}
	canvasContext.Call("stroke", path)
}

func drawHorizLines(canvasContext *js.Value, ycoords []float64, xlow, xhigh float64, color *Color) {
	path := path2dCtor.New()
	canvasContext.Set("strokeStyle", color.toString())
	for _, ycoord := range ycoords {
		addLine2Path(&path, xlow, ycoord+0.5, xhigh, ycoord+0.5)
	}
	canvasContext.Call("stroke", path)
}

func drawLine(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, color *Color) {
	path := path2dCtor.New()
	canvasContext.Set("strokeStyle", color.toString())
	addLine2Path(&path, xlow, ylow, xhigh, yhigh)
	canvasContext.Call("stroke", path)
}

func drawHorizLine(canvasContext *js.Value, xlow, xhigh, y float64, color *Color) {
	path := path2dCtor.New()
	canvasContext.Set("strokeStyle", color.toString())
	addLine2Path(&path, xlow, y+0.5, xhigh, y+0.5)
	canvasContext.Call("stroke", path)
}

func drawVertLine(canvasContext *js.Value, ylow, yhigh, x float64, color *Color) {
	path := path2dCtor.New()
	canvasContext.Set("strokeStyle", color.toString())
	addLine2Path(&path, x+0.5, ylow, x+0.5, yhigh)
	canvasContext.Call("stroke", path)
}

// Adds a line to the path, no actual drawing takes place
func addLine2Path(path *js.Value, xlow, ylow, xhigh, yhigh float64) {
	path.Call("moveTo", xlow, ylow)
	path.Call("lineTo", xhigh, yhigh)
}

func noStrokeFillRect(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, fillColor *Color) {
	path := path2dCtor.New()
	canvasContext.Set("fillStyle", fillColor.toString())
	addRect2Path(&path, xlow, ylow, xhigh, yhigh)
	canvasContext.Call("fill", path)
}

func noStrokeFillRectNoAdjust(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, fillColor *Color) {
	path := path2dCtor.New()
	canvasContext.Set("fillStyle", fillColor.toString())
	addRect2PathNoAdjust(&path, xlow, ylow, xhigh, yhigh)
	canvasContext.Call("fill", path)
}

func strokeFillRect(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, strokeColor, fillColor *Color) {
	path := path2dCtor.New()
	canvasContext.Set("strokeStyle", strokeColor.toString())
	canvasContext.Set("fillStyle", fillColor.toString())
	addRect2Path(&path, xlow, ylow, xhigh, yhigh)
	canvasContext.Call("stroke", path)
	canvasContext.Call("fill", path)
}

func strokeNoFillRect(canvasContext *js.Value, xlow, ylow, xhigh, yhigh float64, strokeColor *Color) {
	path := path2dCtor.New()
	canvasContext.Set("strokeStyle", strokeColor.toString())
	addRect2Path(&path, xlow, ylow, xhigh, yhigh)
	canvasContext.Call("stroke", path)
}

// Adds a rect to the path, no actual drawing takes place
// Adjustment to xlow, ylow to match the adjustments we do with line drawings
func addRect2Path(path *js.Value, xlow, ylow, xhigh, yhigh float64) {
	path.Call("rect", xlow+0.5, ylow+0.5, xhigh-xlow, yhigh-ylow)
}

func addRect2PathNoAdjust(path *js.Value, xlow, ylow, xhigh, yhigh float64) {
	path.Call("rect", xlow, ylow, xhigh-xlow, yhigh-ylow)
}

func drawText(canvasContext *js.Value, xlow, ylow, xhigh, yhigh, xmax, ymax float64,
	text string,
	align TextAlignType, // Used only if attribs is nil
	attribs *CellAttribs) {

	xend, yend := math.Min(xhigh, xmax), math.Min(yhigh, ymax)
	canvasContext.Call("save")
	path := path2dCtor.New()
	addRect2Path(&path, xlow, ylow, xend-1.0, yend-1.0)
	canvasContext.Call("clip", path)

	canvasContext.Set("textBaseline", "bottom")
	startx, starty := xlow, yhigh // yhigh assuming English like language.

	fgcolor := defaultColors.cellStroke
	var bgcolor *Color

	if attribs != nil {
		align = attribs.GetAlignment()
		// Override default FG color
		if fgcolorOverride := attribs.GetFGColor(); fgcolorOverride != nil {
			fgcolor = fgcolorOverride
		}

		bgcolor = attribs.GetBGColor()
	}

	if align == AlignLeft {
		canvasContext.Set("textAlign", "left")
	}
	if align == AlignCenter {
		startx = 0.5 * (xlow + xhigh)
		canvasContext.Set("textAlign", "center")
	} else if align == AlignRight {
		startx = xhigh
		canvasContext.Set("textAlign", "right")
	}

	fontcss := "serif"
	fontSize := constCellFontSizePx
	shouldUnderline := false

	if attribs != nil {
		prefix := ""
		if attribs.IsItalics() {
			prefix = "italic "
		}
		if attribs.IsBold() {
			prefix += "bold "
		}

		shouldUnderline = attribs.IsUnderline()
		fontSize = attribs.GetFontSize()

		if len(prefix) > 0 {
			fontcss = fmt.Sprintf("%s %dpx %s", prefix, fontSize, fontcss)
		} else {
			fontcss = fmt.Sprintf("%dpx %s", fontSize, fontcss)
		}
	} else {
		fontcss = fmt.Sprintf("%dpx %s", fontSize, fontcss)
	}

	if bgcolor != nil {
		noStrokeFillRect(canvasContext, xlow, ylow, xend, yend, bgcolor)
	}

	setFillColor(canvasContext, fgcolor) // For font foreground and underline
	setFont(canvasContext, fontcss)
	canvasContext.Call("fillText", text, startx, starty)
	if shouldUnderline {
		textWidth := canvasContext.Call("measureText", text).Get("width").Float()
		//fmt.Printf("txtwidth = %.1f ", textWidth)
		uxlow := startx
		if align == AlignCenter {
			uxlow -= (textWidth / 2.0)
		} else if align == AlignRight {
			uxlow -= textWidth
		}
		canvasContext.Call("fillRect", uxlow, starty-3, textWidth, 2)
	}
	// Kill the clip path
	canvasContext.Call("restore")

	// DEBUG code
	//strokeNoFillRect(canvasContext, xlow, ylow, xhigh, yhigh, "#0000ff")
	//strokeNoFillRect(canvasContext, xlow, ylow, xmax, ymax, "#ff0000")
}

func setFont(canvasContext *js.Value, fontCSS string) {
	canvasContext.Set("font", fontCSS)
}

func setFillColor(canvasContext *js.Value, fillColor *Color) {
	canvasContext.Set("fillStyle", fillColor.toString())
}

func setStrokeColor(canvasContext *js.Value, strokeColor *Color) {
	canvasContext.Set("strokeStyle", strokeColor.toString())
}

func setLineWidth(canvasContext *js.Value, width float64) {
	canvasContext.Set("lineWidth", width)
}
