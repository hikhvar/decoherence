package prospector

import "strings"

func NewMultiError(errs []error) *MultiError {
	if len(errs) > 0 {
		return &MultiError{
			errors: errs,
		}
	}
	return nil
}

type MultiError struct {
	errors []error
}

func (m *MultiError) Error() string {
	var errs []string
	for _, err := range m.errors {
		errs = append(errs, err.Error())
	}
	return strings.Join(errs, ";")
}
