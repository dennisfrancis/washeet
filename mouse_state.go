package washeet

func (mouse *MouseState) isLeftDown() bool {
	return 0x01 == (*mouse & 0x01)
}

func (mouse *MouseState) isRightDown() bool {
	return 0x02 == (*mouse & 0x02)
}

func (mouse *MouseState) setLeftDown() {
	(*mouse) |= 0x01
}

func (mouse *MouseState) setLeftUp() {
	(*mouse) &= 0xFE
}

func (mouse *MouseState) setRightDown() {
	(*mouse) |= 0x02
}

func (mouse *MouseState) setRightUp() {
	(*mouse) &= 0xFD
}

func defaultMouseState() MouseState {
	return MouseState(0x00)
}
