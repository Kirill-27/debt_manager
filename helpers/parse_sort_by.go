package helpers

import (
	"fmt"
)

func ParseSortBy(sortBy []string) string {
	result := fmt.Sprintf(" ORDER BY ")
	for index, value := range sortBy {
		if index > 0 {
			result += ", "
		}
		if value[0] == '-' {
			result += value[1:] + " DESC"
		} else {
			result += value
		}
	}

	return result
}
