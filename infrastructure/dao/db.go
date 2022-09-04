package dao

import "strings"

func isDuplicateEntryError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate")
}
