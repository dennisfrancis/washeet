package washeet

// NewColor returns a pointer to a new Color object with "red", "blue" and "green" components.
func NewColor(red, green, blue byte) *Color {
	return &Color{
		red:   red,
		green: green,
		blue:  blue,
		alpha: 255,
	}
}

func newRGBAColor(red, green, blue, alpha byte) *Color {
	return &Color{
		red:   red,
		green: green,
		blue:  blue,
		alpha: alpha,
	}
}

var (
	defaultColors = colorSettings{
		gridLine:        "rgba(200, 200, 200, 1.0)",
		cellFill:        "rgba(255, 255, 255, 1.0)",
		cellStroke:      "rgba(0, 0, 0, 1.0)",
		cursorStroke:    "rgba(0, 0, 0, 1.0)",
		headerFill:      "rgba(240, 240, 240, 1.0)",
		selectionStroke: "rgba(100, 156, 244, 1.0)",
		selectionFill:   "rgba(100, 156, 244, 0.2)",
	}
)
