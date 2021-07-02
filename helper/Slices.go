package helper

//SliceStringContains checks if a slice of strings contains a given needle
func SliceStringContains(needle string, haystack []string) bool {
	return SliceStringLinearSearch(needle, haystack) != -1
}

//SliceStringRemove removes all occurrences of the given value from an slice of
//strings
func SliceStringRemove(value string, slice []string) []string {
	var out []string
	for i := 0; i < len(slice); i++ {
		if slice[i] != value {
			out = append(out, slice[i])
		}
	}
	return out
}

//SliceStringLinearSearch finds the first occurrence of the given needle in the
//haystack.
func SliceStringLinearSearch(needle string, haystack []string) int {
	for i := 0; i < len(haystack); i++ {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}

//SliceStringUnique removes all duplicated items in the input slice
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
