package washeet

import "fmt"

// NewColor returns a pointer to a new Color object with "red", "blue" and "green" components.
func NewColor(red, green, blue uint8) *Color {
	color := Color(uint32(red)<<24 | uint32(green)<<16 | uint32(blue)<<8 | uint32(255))
	return &color
}

func newRGBAColor(red, green, blue, alpha uint8) *Color {
	color := Color(uint32(red)<<24 | uint32(green)<<16 | uint32(blue)<<8 | uint32(alpha))
	return &color
}

// Red returns the red component of Color.
func (color *Color) Red() uint8 {
	return uint8(*color >> 24)
}

// Green returns the green component of Color.
func (color *Color) Green() uint8 {
	return uint8((*color >> 16) & Color(0xff))
}

// Blue returns the blue component of Color.
func (color *Color) Blue() uint8 {
	return uint8((*color >> 8) & Color(0xff))
}

// alpha returns the alpha component of Color.
func (color *Color) alpha() float32 {
	return float32(*color&Color(0xff)) / float32(255.0)
}

func (color *Color) toString() string {
	return fmt.Sprintf(
		"rgba(%d, %d, %d, %.1f)",
		color.Red(),
		color.Green(),
		color.Blue(),
		color.alpha(),
	)
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
