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
