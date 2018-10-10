package washeet

func ceilInt64(num int64, den int64) int64 {
	result := num / den
	rem := num % den
	if rem > 0 {
		result++
	}
	return result
}

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

}

func row2RowLabel(nRow int64) string {

}
