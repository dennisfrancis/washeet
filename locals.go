package washeet

import (
	"fmt"
)

var (
	letters = [...]byte{'A', 'B', 'C', 'D', 'E', 'F', 'G',
		'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q',
		'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z'}
)

func minInt64(n1 int64, n2 int64) int64 {
	result := n1
	if n2 < n1 {
		result = n2
	}
	return result
}

func maxInt64(n1 int64, n2 int64) int64 {
	result := n1
	if n2 > n1 {
		result = n2
	}
	return result
}

func col2ColLabel(nCol int64) string {
	label := string(letters[nCol%26])
	for {
		nCol /= 26
		if nCol == 0 {
			break
		}
		if nCol%26 == 0 {
			label = string(letters[25]) + label
			nCol -= 26
		} else {
			label = string(letters[(nCol%26)-1]) + label
		}
	}
	return label
}

func row2RowLabel(nRow int64) string {
	return fmt.Sprintf("%d", nRow+1)
}

// Assumes intervalArray is sorted of course.
func getIntervalIndex(val float64, intervalArray []float64) (index int64) {

	// Result for out of bound cases
	index = -1

	// There must be at least one interval, so there must be at least 2 elements !
	if len(intervalArray) < 2 {
		return
	}

	// Note : intervalArray has len(intervalArray) - 1 intervals.
	lowIntervalIdx, highIntervalIdx := int64(0), int64(len(intervalArray)-2)

	// intervalArray[lowIntervalIdx] is the first interval's start point
	// intervalArray[highIntervalIdx+1] is the last interval's end point
	if val <= intervalArray[lowIntervalIdx] || intervalArray[highIntervalIdx+1] < val {
		return
	}

	// Do binary search to find the right interval
	// TODO: do interpolation search instead !
	for lowIntervalIdx <= highIntervalIdx {
		index = (lowIntervalIdx + highIntervalIdx) / 2
		thisIntervalStart := intervalArray[index]
		thisIntervalEnd := intervalArray[index+1]
		if thisIntervalStart < val && val <= thisIntervalEnd {
			// Found the interval where val lies
			return
		}
		if val <= thisIntervalStart {
			highIntervalIdx = index - 1
		} else {
			lowIntervalIdx = index + 1
		}
	}

	index = -1
	return
}
