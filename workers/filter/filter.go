package filter

import "slices"

type Filter struct {
	DataDuplicate map[uint][]uint
}

func NewFilter() *Filter {
	return &Filter{
		DataDuplicate: make(map[uint][]uint),
	}
}

func (f *Filter) Ok(commentID, userID uint) bool {
	userIDs, ok := f.DataDuplicate[commentID]
	if ok && slices.Contains(userIDs, userID) {
		return false
	}
	if !ok {
		f.DataDuplicate[commentID] = []uint{userID}
	}
	return true
}
