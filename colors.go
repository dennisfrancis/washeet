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
		gridLine:        newRGBAColor(200, 200, 200, 255), // "rgba(200, 200, 200, 1.0)"
		cellFill:        newRGBAColor(255, 255, 255, 255), // "rgba(255, 255, 255, 1.0)"
		cellStroke:      newRGBAColor(0, 0, 0, 255),       // "rgba(0, 0, 0, 1.0)"
		cursorStroke:    newRGBAColor(0, 0, 0, 255),       // "rgba(0, 0, 0, 1.0)"
		headerFill:      newRGBAColor(240, 240, 240, 255), // "rgba(240, 240, 240, 1.0)"
		selectionStroke: newRGBAColor(100, 156, 244, 255), // "rgba(100, 156, 244, 1.0)"
		selectionFill:   newRGBAColor(100, 156, 244, 51),  // "rgba(100, 156, 244, 0.2)"
	}
)
