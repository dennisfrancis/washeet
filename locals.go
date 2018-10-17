package washeet

import (
	"fmt"
)

const (
	letters = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
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
	label := ""
	for {
		label = label + letters[nCol%26]
		nCol /= 26
	}
	return "COL"
}

func row2RowLabel(nRow int64) string {
	return fmt.Sprintf("%d", nRow+1)
}
