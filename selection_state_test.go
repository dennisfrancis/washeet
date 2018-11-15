package washeet

import (
	"testing"
)

func TestDefaultSelectionState(t *testing.T) {

	state := defaultSelectionState()
	start := state.getRefStartCell()
	curr := state.getRefCurrCell()
	if start.col != 0 || start.row != 0 || curr.col != 0 || curr.row != 0 {
		t.Error("Wrong default selection state : ", state)
	}
}

func TestSetRefStartCurrCell(t *testing.T) {

	cases := []cellCoords{
		cellCoords{0, 10},
		cellCoords{10, 0},
		cellCoords{5, 10},
		cellCoords{10, 5},
	}
	state := defaultSelectionState()
	start := state.getRefStartCell()
	curr := state.getRefCurrCell()

	for _, casStart := range cases {

		state.setRefStartCell(casStart.col, casStart.row)

		for _, casCurr := range cases {

			state.setRefCurrCell(casCurr.col, casCurr.row)

			if start.col != casStart.col || start.row != casStart.row {
				t.Errorf("Wrong start cell state : expected (%d, %d) but got (%d, %d)", casStart.col, casStart.row, start.col, start.row)
			}

			if curr.col != casCurr.col || curr.row != casCurr.row {
				t.Errorf("Wrong curr cell state : expected (%d, %d) but got (%d, %d)", casCurr.col, casCurr.row, curr.col, curr.row)
			}
		}
	}
}
