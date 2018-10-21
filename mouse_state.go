package washeet

func (self *MouseState) setLeftDown() {
	*self |= 0x01
}

func (self *MouseState) setLeftUp() {
	*self &= 0xFE
}

func (self *MouseState) setRightDown() {
	*self |= 0x02
}

func (self *MouseState) setRightUp() {
	*self &= 0xFD
}

func defaultMouseState() MouseState {
	return MouseState(0x00)
}
