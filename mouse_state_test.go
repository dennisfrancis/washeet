package washeet

import "testing"

func TestMouseDefaultState(t *testing.T) {

	state := defaultMouseState()
	if state.isLeftDown() || state.isRightDown() {
		t.Error("Wrong default mouse state : ", state)
	}
}

func TestMouseLeftButtonState(t *testing.T) {

	state := defaultMouseState()
	state.setLeftDown()
	state.setRightUp()
	if !state.isLeftDown() {
		t.Error("Left button not down")
	}

	state.setRightDown()
	if !state.isLeftDown() {
		t.Error("Left button not down")
	}

	state.setLeftUp()
	state.setRightDown()
	if state.isLeftDown() {
		t.Error("Left button not up")
	}

	state.setRightUp()
	if state.isLeftDown() {
		t.Error("Left button not up")
	}
}

func TestMouseRightButtonState(t *testing.T) {

	state := defaultMouseState()
	state.setRightDown()
	state.setLeftUp()
	if !state.isRightDown() {
		t.Error("Right button not down")
	}

	state.setLeftDown()
	if !state.isRightDown() {
		t.Error("Right button not down")
	}

	state.setRightUp()
	state.setLeftDown()
	if state.isRightDown() {
		t.Error("Right button not up")
	}

	state.setLeftUp()
	if state.isRightDown() {
		t.Error("Right button not up")
	}
}
