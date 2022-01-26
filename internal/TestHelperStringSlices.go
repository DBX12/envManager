package internal

import (
	"testing"
)

//AssertStringSliceEqual tests if all elements in slice "want" are present in
//slice "got" and vice versa. It will not pay attention to the order of the
//elements; this is intentional. If you want to check that both slices are in
//the same order, use reflect.DeepEqual.
//Will return a boolean indicating if the slices are equal.
//The logic is a copy from helper.SliceStringEqual
func AssertStringSliceEqual(t *testing.T, want []string, got []string) bool {
	t.Helper()
	lenA := len(want)
	assertionIsValid := true
	if lenA != len(got) {
		t.Error("Slices do not have the same length")
		return false
	}
	for i := 0; i < lenA; i++ {
		if sliceStringLinearSearch(t, want[i], got) == -1 {
			// entry want[i] is missing from got, the slices are not equal
			t.Errorf("String %s is present in the 'want' slice but missing from the 'got' slice", want[i])
			assertionIsValid = false
		}
	}
	return assertionIsValid
}

//This is a copy from helper.SliceStringLinearSearch with the addition of the
//parameter t *testing.T
func sliceStringLinearSearch(t *testing.T, needle string, haystack []string) int {
	t.Helper()
	for i := 0; i < len(haystack); i++ {
		if haystack[i] == needle {
			return i
		}
	}
	return -1
}
