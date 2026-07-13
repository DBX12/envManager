package helper

import (
	"slices"

	"github.com/manifoldco/promptui"
)

// Multiselect uses promptui.Select to let the user select one or more entries
// from a list of options. If the user selects no entry (immediately selects
// "DONE"), you will receive an empty slice and no error.
func Multiselect(options []string, prompt string) ([]string, error) {
	const DoneValue = "DONE"

	prompt = prompt + " (select " + DoneValue + " when you are done)"
	options = append([]string{DoneValue}, options...)

	var selections []string
	var sel string
	for {
		var err error
		// rebuild the prompt on every iteration because we remove selected
		// values from the options slice.
		p := promptui.Select{
			Label: prompt,
			Items: options,
		}
		_, sel, err = p.Run()
		if err != nil {
			return nil, err
		}

		if sel == DoneValue {
			break
		}
		selections = append(selections, sel)
		idx := slices.Index(options, sel)
		options = slices.Delete(options, idx, idx+1)
	}
	return selections, nil
}
