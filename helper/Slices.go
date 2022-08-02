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
//haystack. Returns -1 if the needle was not found.
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

//SliceStringEqual checks if the slices contain the same string values without
//considering the order of the items. This means [foo bar] is equal to
//[bar foo] but not equal to [foo baz].
func SliceStringEqual(a []string, b []string) bool {
	lenA := len(a)
	if lenA != len(b) {
		return false
	}
	for i := 0; i < lenA; i++ {
		if SliceStringLinearSearch(a[i], b) == -1 {
			// entry a[i] is missing from b, the slices are not equal
			return false
		}
	}
	return true
}

//SliceStringReverse reverses the input slice. If the input slice is nil, an empty slice is returned.
func SliceStringReverse(input []string) []string {
	length := len(input)
	output := make([]string, length)
	for i := 0; i < length; i++ {
		output[length-i-1] = input[i]
	}
	return output
}
