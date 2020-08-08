package simulator

import (
	"os/exec"
	"strconv"
)

// Error is returned by LookPath when it fails to classify a file as an
// executable.
type Error struct {
	// Name is the file name for which the error occurred.
	Name string
	// Err is the underlying error.
	Err error
}

func (e *Error) Error() string {
	return "exec: " + strconv.Quote(e.Name) + ": " + e.Err.Error()
}

func (e *Error) Unwrap() error { return e.Err }

// NewCommand todo
func NewCommand(name string, arg ...string) (*exec.Cmd, error) {
	cmd := &exec.Cmd{}
	// TODO
	return cmd, nil
}
