package cmdutil

import "fmt"

// ExitError represents a command failure with an exit code.
type ExitError struct {
	Code int
	Err  error
}

func (e *ExitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return fmt.Sprintf("exit code %d", e.Code)
}

func (e *ExitError) Unwrap() error { return e.Err }

// NewExitError creates a new exit error.
func NewExitError(code int, msg string) *ExitError {
	return &ExitError{Code: code, Err: fmt.Errorf("%s", msg)}
}
