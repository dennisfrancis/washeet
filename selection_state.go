package washeet

func (self *SelectionState) setRefStartCell(col, row int64) {
	self.refStartCell.Col, self.refStartCell.Row = col, row
}

func (self *SelectionState) setRefCurrCell(col, row int64) {
	self.refCurrCell.Col, self.refCurrCell.Row = col, row
}

func (self *SelectionState) getRefStartCell() *CellCoords {
	return &(self.refStartCell)
}

func (self *SelectionState) getRefCurrCell() *CellCoords {
	return &(self.refCurrCell)
}

func defaultSelectionState() SelectionState {
	return SelectionState{refStartCell: CellCoords{0, 0}, refCurrCell: CellCoords{0, 0}}
}
