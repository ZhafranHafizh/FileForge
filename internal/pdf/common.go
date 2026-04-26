package pdf

import (
	"errors"
)

type ValidationError struct {
	Err error
}

func (e ValidationError) Error() string {
	return e.Err.Error()
}

func (e ValidationError) Unwrap() error {
	return e.Err
}

type DependencyError struct {
	Err error
}

func (e DependencyError) Error() string {
	return e.Err.Error()
}

func (e DependencyError) Unwrap() error {
	return e.Err
}

func IsValidationError(err error) bool {
	var target ValidationError
	return errors.As(err, &target)
}

func IsDependencyError(err error) bool {
	var target DependencyError
	return errors.As(err, &target)
}
