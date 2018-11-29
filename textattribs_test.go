package washeet

import (
	"testing"
)

func TestDefaultTextAttributes(t *testing.T) {
	tattribs := newDefaultTextAttributes()
	if tattribs.isBold() {
		t.Error("should not have bold set")
	}

	if tattribs.isItalics() {
		t.Error("should not have italics set")
	}

	if tattribs.isUnderline() {
		t.Error("should not have underline set")
	}
}
