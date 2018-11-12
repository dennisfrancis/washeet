package washeet

func (selection *SelectionState) setRefStartCell(col, row int64) {
	selection.refStartCell.Col, selection.refStartCell.Row = col, row
}

func (selection *SelectionState) setRefCurrCell(col, row int64) {
	selection.refCurrCell.Col, selection.refCurrCell.Row = col, row
}

func (selection *SelectionState) getRefStartCell() *CellCoords {
	return &(selection.refStartCell)
}

func (selection *SelectionState) getRefCurrCell() *CellCoords {
	return &(selection.refCurrCell)
}

func defaultSelectionState() SelectionState {
	return SelectionState{refStartCell: CellCoords{0, 0}, refCurrCell: CellCoords{0, 0}}
}
