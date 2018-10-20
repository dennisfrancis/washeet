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

func TestGetIntervalIndex(t *testing.T) {
	type queryResultPair struct {
		value float64
		index int64
	}
	type intervalSearchTestCase struct {
		interval []float64
		subcases []queryResultPair
	}

	cases := []intervalSearchTestCase{

		// Bad interval array
		intervalSearchTestCase{
			[]float64{},
			[]queryResultPair{
				queryResultPair{10.0, -1},
				queryResultPair{-100.0, -1},
				queryResultPair{1000.0, -1},
			},
		},

		// Bad interval array
		intervalSearchTestCase{
			[]float64{20.0},
			[]queryResultPair{
				queryResultPair{10.0, -1},
				queryResultPair{-100.0, -1},
				queryResultPair{1000.0, -1},
			},
		},

		// Valid interval array
		intervalSearchTestCase{
			[]float64{10.2, 11.5, 100.57, 2000.5, 2002.67},
			[]queryResultPair{
				// Out of bound cases
				queryResultPair{-20.3, -1},
				queryResultPair{5.0, -1},
				queryResultPair{10.2, -1}, // start element of interval array is not inclusive
				queryResultPair{2002.679, -1},
				queryResultPair{4000.5, -1},
				// Non out of bound cases
				queryResultPair{10.21, 0},
				queryResultPair{10.9, 0},
				queryResultPair{11.5, 0}, // end element of the interval is inclusive
				queryResultPair{11.51, 1},
				queryResultPair{50.0, 1},
				queryResultPair{100.57, 1},
				queryResultPair{1000.49, 2},
				queryResultPair{2000.5, 2},
				queryResultPair{2000.51, 3},
				queryResultPair{2001.54, 3},
				queryResultPair{2002.67, 3},
			},
		},
	}

	for _, cas := range cases {
		for _, pair := range cas.subcases {
			value, index := pair.value, pair.index
			result := getIntervalIndex(value, cas.interval)
			if result != index {
				t.Error("For search point", value, "expected interval index =", index, "but got", result)
			}
		}
	}
}
