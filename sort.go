package cr

import (
	"fmt"
	"strings"
)

type SortBy []Sort

type Sort struct {
	Column string
	Order  string
}

func (s Sort) String() string {
	if s.Order == "" {
		return fmt.Sprintf("%s ASC", s.Column)
	}
	return fmt.Sprintf("%s %s", s.Column, s.Order)
}

func (by SortBy) Strs() []string {
	return apply(by, func(item Sort) string {
		return item.String()
	})
}

func (by SortBy) String() string {
	return strings.Join(by.Strs(), ", ")
}

func ParseSort(sort string) (s []Sort) {
	return apply(strings.Split(sort, ","), func(column string) Sort {
		column = strings.TrimSpace(column)
		if column[0] == '-' {
			return Sort{Column: column[1:], Order: "DESC"}
		}
		return Sort{Column: column, Order: "ASC"}
	})
}

func apply[T any, R any](collection []T, iteratee func(item T) R) []R {
	if len(collection) == 0 {
		return nil
	}

	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item)
	}

	return result
}
