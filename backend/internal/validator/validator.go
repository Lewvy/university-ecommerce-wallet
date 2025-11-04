package validator

import (
	"regexp"
	"slices"
)

var (
	EmailRX      = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	NonNumericRX = regexp.MustCompile(`[^\d]`)
)

type ValidationError struct {
	Errors map[string]string
}

func (v *ValidationError) Error() string {
	return "validation error"
}

func New() *ValidationError {
	return &ValidationError{Errors: make(map[string]string)}
}

func (v *ValidationError) Valid() bool {
	return len(v.Errors) == 0
}

func (v *ValidationError) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *ValidationError) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func (v *ValidationError) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
