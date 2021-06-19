package helper

import (
	"github.com/manifoldco/promptui"
)

//input is an abstraction layer for retrieving user input. It is to be used as
//singleton via GetInput.
type input struct {
	Inputs []string
}

//hasPresetInputValues checks that the Inputs slice is defined and contains at
//least one entry.
func (i *input) hasPresetInputValues() bool {
	return i.Inputs != nil && len(i.Inputs) > 0
}

//getPresetInputValue returns the first value in the Inputs slice.
//This method will panic if the slice is not defined or empty!
//Check with hasPresetInputValues that there are presets before calling this method.
func (i *input) getPresetInputValue() string {
	out := i.Inputs[0]
	if len(i.Inputs) >= 2 {
		// there is at least one other item in the slice, so just pop the first entry
		i.Inputs = i.Inputs[1:]
	} else {
		// this was the last preset, let the GC take the slice
		i.Inputs = nil
	}
	return out
}

//PromptPassword prompts for a password and hides the input with the mask value.
//Setting the mask to 0 disables the masking. If you are doing that, you can
//call PromptString as well.
func (i *input) PromptPassword(prompt string, mask rune) (string, error) {
	if i.hasPresetInputValues() {
		return i.getPresetInputValue(), nil
	}
	promptUi := promptui.Prompt{
		Label: prompt,
		Mask:  mask,
	}

	result, err := promptUi.Run()

	if err != nil {
		return "", err
	}
	return result, nil
}

//PromptString prompts the user for an input
func (i *input) PromptString(prompt string) (string, error) {
	return i.PromptPassword(prompt, 0)
}
