package washeet

import "testing"

func TestCol2ColLabel(t *testing.T) {
	type colLabelCases struct {
		col       int64
		colString string
		label     string
	}
	cases := []colLabelCases{
		colLabelCases{
			0 + 1*26,
			"0 + 1*26",
			"AA",
		},
		colLabelCases{
			25 + 1*26,
			"25 + 1*26",
			"AZ",
		},
		colLabelCases{
			25 + 26*26,
			"25 + 26*26",
			"ZZ",
		},
		colLabelCases{
			0 + 1*26 + 1*26*26,
			"0 + 1*26 + 1*26*26",
			"AAA",
		},
		colLabelCases{
			25 + 26*26 + 1*26*26,
			"25 + 26*26 + 1*26*26",
			"AZZ",
		},
		colLabelCases{
			25 + 2*26 + 26*26*26,
			"25 + 2*26 + 26*26*26",
			"ZBZ",
		},
		colLabelCases{
			25 + 2*26 + 3*26*26 + 25*26*26*26,
			"25 + 2*26 + 3*26*26 + 25*26*26*26",
			"YCBZ",
		},
	}

	for _, cas := range cases {
		result := col2ColLabel(cas.col)
		if result != cas.label {
			t.Error("For", cas.colString, "expected", cas.label, "got", result)
		}
	}
}
