package washeet

// NewDefaultCellAttribs return a pointer to a new instance of CellAttribs
// with default settings.
func NewDefaultCellAttribs() *CellAttribs {
	cAttribs := &CellAttribs{
		txtAttribs: newDefaultTextAttributes(),
	}
	return cAttribs
}

// SetBold sets/unsets the cell text's bold attribute.
func (cAttribs *CellAttribs) SetBold(flag bool) {
	cAttribs.txtAttribs.setBold(flag)
}

// SetItalics sets/unsets the cell text's italics attribute.
func (cAttribs *CellAttribs) SetItalics(flag bool) {
	cAttribs.txtAttribs.setItalics(flag)
}

// SetUnderline sets/unsets the cell text's underline attribute.
func (cAttribs *CellAttribs) SetUnderline(flag bool) {
	cAttribs.txtAttribs.setUnderline(flag)
}

// IsBold returns whether bold attribute is set or not.
func (cAttribs *CellAttribs) IsBold() bool {
	return cAttribs.txtAttribs.isBold()
}

// IsItalics returns whether italics attribute is set or not.
func (cAttribs *CellAttribs) IsItalics() bool {
	return cAttribs.txtAttribs.isItalics()
}

// IsUnderline returns whether underline attribute is set or not.
func (cAttribs *CellAttribs) IsUnderline() bool {
	return cAttribs.txtAttribs.isUnderline()
}
