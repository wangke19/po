package cmdutil

import "fmt"

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

func NewExitError(code int, msg string) *ExitError {
	return &ExitError{Code: code, Err: fmt.Errorf("%s", msg)}
}
