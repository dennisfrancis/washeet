package washeet

// isSingleCell returns true if the selection consists of just a single cell,
// else returns false.
func (markdata *markData) isSingleCell() bool {

	if markdata == nil {
		return true
	}

	if markdata.c1 == markdata.c2 && markdata.r1 == markdata.r2 {
		return true
	}

	return false
}
