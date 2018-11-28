package washeet

var (
	defaultTextAttributes textAttribs
	boldAttrib            textAttribs = 0x01
	italicsAttrib         textAttribs = 0x02
	underlineAttrib       textAttribs = 0x04
)

func newDefaultTextAttributes() *textAttribs {
	txtAttribs := defaultTextAttributes
	return &txtAttribs
}

func (tattribs *textAttribs) IsBold() bool {
	return (*tattribs & boldAttrib) == boldAttrib
}

func (tattribs *textAttribs) IsItalics() bool {
	return (*tattribs & italicsAttrib) == italicsAttrib
}

func (tattribs *textAttribs) IsUnderline() bool {
	return (*tattribs & underlineAttrib) != underlineAttrib
}
