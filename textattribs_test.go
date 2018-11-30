package washeet

import (
	"testing"
)

func TestDefaultTextAttributes(t *testing.T) {
	tattribs := newDefaultTextAttributes()
	if tattribs.isBold() {
		t.Error("should not have bold set")
	}

	if tattribs.isItalics() {
		t.Error("should not have italics set")
	}

	if tattribs.isUnderline() {
		t.Error("should not have underline set")
	}
}

func TestSetUnsetAttributes(t *testing.T) {
	type flagSetType uint8
	const (
		dndFlag flagSetType = iota
		setFlag
		unsetFlag
	)

	type testCase struct {
		tattribs        textAttribs
		boldSet         flagSetType
		italicsSet      flagSetType
		underlineSet    flagSetType
		boldExpect      bool
		italicsExpect   bool
		underlineExpect bool
	}

	tcases := []testCase{
		// Test setting of flags one by one when all others are unset.
		// #0
		testCase{
			tattribs:        defaultTextAttributes,
			boldSet:         setFlag,
			italicsSet:      dndFlag,
			underlineSet:    dndFlag,
			boldExpect:      true,
			italicsExpect:   false,
			underlineExpect: false,
		},
		// #1
		testCase{
			tattribs:        defaultTextAttributes,
			boldSet:         dndFlag,
			italicsSet:      setFlag,
			underlineSet:    dndFlag,
			boldExpect:      false,
			italicsExpect:   true,
			underlineExpect: false,
		},
		// #2
		testCase{
			tattribs:        defaultTextAttributes,
			boldSet:         dndFlag,
			italicsSet:      dndFlag,
			underlineSet:    setFlag,
			boldExpect:      false,
			italicsExpect:   false,
			underlineExpect: true,
		},

		// Test getting of flags one by one when all are set.
		// #3
		testCase{
			tattribs:        0x07,
			boldSet:         dndFlag,
			italicsSet:      dndFlag,
			underlineSet:    dndFlag,
			boldExpect:      true,
			italicsExpect:   true,
			underlineExpect: true,
		},

		// Test unsetting of flags one by one when all are set.
		// #4
		testCase{
			tattribs:        0x07,
			boldSet:         unsetFlag,
			italicsSet:      dndFlag,
			underlineSet:    dndFlag,
			boldExpect:      false,
			italicsExpect:   true,
			underlineExpect: true,
		},

		// #5
		testCase{
			tattribs:        0x07,
			boldSet:         dndFlag,
			italicsSet:      unsetFlag,
			underlineSet:    dndFlag,
			boldExpect:      true,
			italicsExpect:   false,
			underlineExpect: true,
		},

		// #6
		testCase{
			tattribs:        0x07,
			boldSet:         dndFlag,
			italicsSet:      dndFlag,
			underlineSet:    unsetFlag,
			boldExpect:      true,
			italicsExpect:   true,
			underlineExpect: false,
		},
	}

	for idx, tcase := range tcases {
		attrib := tcase.tattribs

		switch tcase.boldSet {
		case setFlag:
			attrib.setBold(true)
		case unsetFlag:
			attrib.setBold(false)
		}

		switch tcase.italicsSet {
		case setFlag:
			attrib.setItalics(true)
		case unsetFlag:
			attrib.setItalics(false)
		}

		switch tcase.underlineSet {
		case setFlag:
			attrib.setUnderline(true)
		case unsetFlag:
			attrib.setUnderline(false)
		}

		if attrib.isBold() != tcase.boldExpect {
			t.Errorf("testcase#%d : isBold() : expected = %v, got = %v", idx, tcase.boldExpect, attrib.isBold())
		}

		if attrib.isItalics() != tcase.italicsExpect {
			t.Errorf("testcase#%d : isItalics() : expected = %v, got = %v", idx, tcase.italicsExpect, attrib.isItalics())
		}

		if attrib.isUnderline() != tcase.underlineExpect {
			t.Errorf("testcase#%d : isUnderline() : expected = %v, got = %v", idx, tcase.underlineExpect, attrib.isUnderline())
		}
	}

}
