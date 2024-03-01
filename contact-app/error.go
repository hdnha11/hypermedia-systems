package main

import "fmt"

type ValidationError struct {
	errors map[string]string
}

func (e *ValidationError) Error() string {
	var errMessage string
	for field, message := range e.errors {
		errMessage += fmt.Sprintf("Field %s: %s\n", field, message)
	}

	return errMessage
}

func (e *ValidationError) AddFieldError(field, message string) {
	if e.errors == nil {
		e.errors = make(map[string]string)
	}

	e.errors[field] = message
}

func (e *ValidationError) FieldErrors() map[string]string {
	return e.errors
}

func FieldErrors(err error) map[string]string {
	if verr, ok := err.(*ValidationError); ok {
		return verr.FieldErrors()
	}

	return map[string]string{}
}
