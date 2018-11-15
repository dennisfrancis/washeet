package washeet

func (mouse *mouseState) isLeftDown() bool {
	return 0x01 == (*mouse & 0x01)
}

func (mouse *mouseState) isRightDown() bool {
	return 0x02 == (*mouse & 0x02)
}

func (mouse *mouseState) setLeftDown() {
	(*mouse) |= 0x01
}

func (mouse *mouseState) setLeftUp() {
	(*mouse) &= 0xFE
}

func (mouse *mouseState) setRightDown() {
	(*mouse) |= 0x02
}

func (mouse *mouseState) setRightUp() {
	(*mouse) &= 0xFD
}

func defaultMouseState() mouseState {
	return mouseState(0x00)
}
