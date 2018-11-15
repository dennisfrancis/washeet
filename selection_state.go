package washeet

func (selection *selectionState) setRefStartCell(col, row int64) {
	selection.refStartCell.col, selection.refStartCell.row = col, row
}

func (selection *selectionState) setRefCurrCell(col, row int64) {
	selection.refCurrCell.col, selection.refCurrCell.row = col, row
}

func (selection *selectionState) getRefStartCell() *cellCoords {
	return &(selection.refStartCell)
}

func (selection *selectionState) getRefCurrCell() *cellCoords {
	return &(selection.refCurrCell)
}

func defaultSelectionState() selectionState {
	return selectionState{refStartCell: cellCoords{0, 0}, refCurrCell: cellCoords{0, 0}}
}
