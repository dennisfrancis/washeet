package washeet

func (self *MarkData) IsSingleCell() bool {

	if self == nil {
		return true
	}

	if self.C1 == self.C2 && self.R1 == self.R2 {
		return true
	}

	return false
}
