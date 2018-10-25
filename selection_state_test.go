package washeet

import (
	"testing"
)

func TestDefaultSelectionState(t *testing.T) {

	state := defaultSelectionState()
	start := state.getRefStartCell()
	curr := state.getRefCurrCell()
	if start.Col != 0 || start.Row != 0 || curr.Col != 0 || curr.Row != 0 {
		t.Error("Wrong default selection state : ", state)
	}
}

func TestSetRefStartCurrCell(t *testing.T) {

	cases := []CellCoords{
		CellCoords{0, 10},
		CellCoords{10, 0},
		CellCoords{5, 10},
		CellCoords{10, 5},
	}
	state := defaultSelectionState()
	start := state.getRefStartCell()
	curr := state.getRefCurrCell()

	for _, casStart := range cases {

		state.setRefStartCell(casStart.Col, casStart.Row)

		for _, casCurr := range cases {

			state.setRefCurrCell(casCurr.Col, casCurr.Row)

			if start.Col != casStart.Col || start.Row != casStart.Row {
				t.Errorf("Wrong start cell state : expected (%d, %d) but got (%d, %d)", casStart.Col, casStart.Row, start.Col, start.Row)
			}

			if curr.Col != casCurr.Col || curr.Row != casCurr.Row {
				t.Errorf("Wrong curr cell state : expected (%d, %d) but got (%d, %d)", casCurr.Col, casCurr.Row, curr.Col, curr.Row)
			}
		}
	}
}
