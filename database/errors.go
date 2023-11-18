package database

import "fmt"

type (
	MissingRequiredOptionError struct {
		optionName string
	}
)

func (e MissingRequiredOptionError) Error() string {
	return fmt.Sprintf("required option not defined: %s", e.optionName)
}

func (e MissingRequiredOptionError) Is(template error) bool {
	if template, ok := template.(MissingRequiredOptionError); ok {
		return e.optionName == "" || e.optionName == template.optionName
	}

	return false
}
