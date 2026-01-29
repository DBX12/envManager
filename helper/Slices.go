package helper

import "slices"

// SliceStringRemove removes all occurrences of the given value from a slice of
// strings
func SliceStringRemove(value string, slice []string) []string {
	return slices.DeleteFunc(slices.Clone(slice), func(s string) bool {
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
