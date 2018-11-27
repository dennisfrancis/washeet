package washeet

import (
	"fmt"
	"testing"
)

func TestColorRGB(t *testing.T) {
	red, green, blue, alpha := uint8(232), uint8(34), uint8(56), uint8(255)
	color := NewColor(red, green, blue)

	if color.Red() != red {
		t.Error("Red() : Expected", red, "got", color.Red())
	}

	if color.Green() != green {
		t.Error("Green() : Expected", green, "got", color.Green())
	}

	if color.Blue() != blue {
		t.Error("Blue() : Expected", blue, "got", color.Blue())
	}

	alphaf := float32(alpha) / float32(255)
	if color.alpha() != alphaf {
		t.Error("alpha() : Expected", alphaf, "got", color.alpha())
	}

	colorString := fmt.Sprintf("rgba(%d, %d, %d, %.1f)", red, green, blue, alphaf)
	if color.toString() != colorString {
		t.Error("toString() : Expected", colorString, color.toString())
	}
}
