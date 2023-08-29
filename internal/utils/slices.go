package utils

import "strings"

func Filter[T any](a []T, predicate func(b T) bool) *T {
	for _, v := range a {
		v := v
		if predicate(v) {
			return &v
		}
	}
	return nil
}

func AnyMatch[T any](a []T, predicate func(b T) bool) bool {
	match := Filter(a, predicate)
	return match != nil
}

func AnyMatchStringCaseInsensitive(a []string, match string) bool {
	return AnyMatch(a, func(b string) bool {
		return strings.EqualFold(b, match)
	})
}
