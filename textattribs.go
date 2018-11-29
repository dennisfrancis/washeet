package washeet

const (
	defaultTextAttributes textAttribs = iota
	boldAttrib            textAttribs = 1 << (iota - 1) // 1
	italicsAttrib                                       // 2
	underlineAttrib                                     // 4
)

func newDefaultTextAttributes() *textAttribs {
	txtAttribs := defaultTextAttributes
	return &txtAttribs
}

func (tattribs *textAttribs) isBold() bool {
	return (*tattribs & boldAttrib) == boldAttrib
}

func (tattribs *textAttribs) isItalics() bool {
	return (*tattribs & italicsAttrib) == italicsAttrib
}

func (tattribs *textAttribs) isUnderline() bool {
	return (*tattribs & underlineAttrib) == underlineAttrib
}

func (tattribs *textAttribs) setField(flag bool, field textAttribs) {
	if flag {
		*tattribs |= field
	} else {
		*tattribs &^= field
	}
}

func (tattribs *textAttribs) setBold(flag bool) {
	tattribs.setField(flag, boldAttrib)
}

func (tattribs *textAttribs) setItalics(flag bool) {
	tattribs.setField(flag, italicsAttrib)
}

func (tattribs *textAttribs) setUnderline(flag bool) {
	tattribs.setField(flag, underlineAttrib)
}
