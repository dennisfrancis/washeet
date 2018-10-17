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
