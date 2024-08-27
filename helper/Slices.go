package helper

import "slices"

// SliceStringContains checks if a slice of strings contains a given needle
func SliceStringContains(needle string, haystack []string) bool {
	return slices.Contains(haystack, needle)
}

// SliceStringRemove removes all occurrences of the given value from an slice of
// strings
func SliceStringRemove(value string, slice []string) []string {
	return slices.DeleteFunc(slice, func(s string) bool {
		return s == value
	})
}

// SliceStringUnique removes all duplicated items in the input slice
func SliceStringUnique(input []string) []string {
	temp := make(map[string]bool, len(input))
	output := make([]string, 0, len(input))
	for _, value := range input {
		_, exists := temp[value]
		if !exists {
			temp[value] = true
			output = append(output, value)
		}
	}
	return output
}

// SliceStringEqual checks if the slices contain the same string values without
// considering the order of the items. This means [foo bar] is equal to
// [bar foo] but not equal to [foo baz].
func SliceStringEqual(a []string, b []string) bool {
	// return early if the lengths do not match to avoid costly sorting and comparing
	if len(a) != len(b) {
		return false
	}
	slices.Sort(a)
	slices.Sort(b)
	return slices.Compare(a, b) == 0
}

// SliceStringReverse reverses the input slice. If the input slice is nil, an empty slice is returned.
func SliceStringReverse(input []string) []string {
	slices.Reverse(input)
	return input
}
