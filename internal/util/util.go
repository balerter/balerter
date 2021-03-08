package util

import (
	"strings"
)

// CheckUnique checks a slice for unique values,
// If founded non unique elements, returns a conflict element name. Else returns empty string
func CheckUnique(data []string) string {
	m := map[string]struct{}{}

	for _, item := range data {
		item = strings.ToLower(item)
		if _, ok := m[item]; ok {
			return item
		}
		m[item] = struct{}{}
	}

	return ""
}

// InArray returns true, if value in arr
func InArray(value string, arr []string) bool {
	for _, v := range arr {
		if value == v {
			return true
		}
	}

	return false
}
