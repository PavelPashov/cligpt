package cligpt

import (
	"errors"
	"log"

	"github.com/manifoldco/promptui"
)

type isValidInputString func(string) bool

type promptInputContent struct {
	errorMsg           string
	label              string
	isValidInputString isValidInputString
}

type promptSelectContent struct {
	label        string
	selectValues []string
}

type promptSelectReturnType struct {
	index int
	value string
}

func promptGetInput(pc promptInputContent) string {
	validate := func(input string) error {
		if len(input) <= 0 {
			return errors.New(pc.errorMsg)
		}
		if pc.isValidInputString != nil && !pc.isValidInputString(input) {
			return errors.New(pc.errorMsg)
		}
		return nil
	}

	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     pc.label,
		Templates: templates,
		Validate:  validate,
	}

	result, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	return result
}

func promptGetSelect(pc promptSelectContent) promptSelectReturnType {
	index := 0
	var err error
	var result string

	prompt := promptui.Select{
		Label: pc.label,
		Items: pc.selectValues,
	}

	index, result, err = prompt.Run()

	if err != nil {
		log.Fatal(err)
	}

	return promptSelectReturnType{index: index, value: result}
}
