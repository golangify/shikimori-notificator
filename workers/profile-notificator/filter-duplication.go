package profilenotificator

import "slices"

type DuplicatorFilter struct {
	Data map[uint][]uint
}

func NewDuplicatorFilter() *DuplicatorFilter {
	return &DuplicatorFilter{
		Data: make(map[uint][]uint),
	}
}

func (f *DuplicatorFilter) Check(commentID, userID uint) bool {
	userIDs, ok := f.Data[commentID]
	if ok && slices.Contains(userIDs, userID) {
		return true
	}
	if ok {
		userIDs = append(userIDs, userID)
		f.Data[commentID] = userIDs
	} else {
		f.Data[commentID] = []uint{userID}
	}
	return false
}
