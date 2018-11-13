package washeet

func (markdata *MarkData) IsSingleCell() bool {

	if markdata == nil {
		return true
	}

	if markdata.C1 == markdata.C2 && markdata.R1 == markdata.R2 {
		return true
	}

	return false
}
