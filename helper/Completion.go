package helper

import "strings"

//Completion filters possibleValues by removing excludedValues and values not
//matching withPrefix. excludedValues can be set to nil if no values are to be
//excluded. withPrefix can be set to empty string do disable filtering by prefix
func Completion(possibleValues []string, excludedValues []string, withPrefix string) []string {
	var out []string
	for _, value := range possibleValues {
		if excludedValues != nil && SliceStringContains(value, excludedValues) {
			// skip this value as it is excluded
			continue
		}
		if withPrefix != "" && !strings.HasPrefix(value, withPrefix) {
			// we are filtering with a prefix and the current value does not have it
			continue
		}
		out = append(out, value)
	}
	return out
}
