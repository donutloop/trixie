package trixie

import "fmt"

// BadPathError creates error for bad path
type BadPathError struct {
	s string
}

func (bme *BadPathError) Error() string { return fmt.Sprintf("Path is invaild (%s)", bme.s) }

// NewBadPathError returns an error that formats as the given text.
func NewBadPathError(text string) error {
	return &BadPathError{s: text}
}

// BadMethodError creates error for bad method
type BadMethodError struct{}

func (bme *BadMethodError) Error() string { return fmt.Sprintf("Method not vaild") }

// NewBadMethodError returns an error that formats as the given text.
func NewBadMethodError() error {
	return &BadMethodError{}
}
