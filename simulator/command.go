package simulator

import (
	"context"
	"os/exec"
	"path/filepath"
	"strconv"
)

// Error is returned by LookPath when it fails to classify a file as an
// executable.
type Error struct {
	// Name is the file name for which the error occurred.
	Name string
	// Err is the unsmrlying error.
	Err error
}

func (e *Error) Error() string {
	return "exec: " + strconv.Quote(e.Name) + ": " + e.Err.Error()
}

func (e *Error) Unwrap() error { return e.Err }

// NewCommand todo
func (sm *Simulator) NewCommand(name string, arg ...string) (*exec.Cmd, error) {
	cmd := &exec.Cmd{
		Path: name,
		Args: append([]string{name}, arg...),
		Env:  sm.Environ(),
	}
	if filepath.Base(name) != name {
		return cmd, nil
	}
	if file, err := sm.LookPath(name); err == nil {
		cmd.Path = file
	}
	return cmd, nil
}

// NewCommandContext todo
func (sm *Simulator) NewCommandContext(ctx context.Context, name string, arg ...string) (*exec.Cmd, error) {
	file := name
	var err error
	if filepath.Base(name) == name {
		if file, err = sm.LookPath(name); err != nil {
			return nil, exec.ErrNotFound
		}
	}
	cmd := exec.CommandContext(ctx, file, arg...)
	cmd.Env = sm.Environ()
	cmd.Args[0] = name // reset arg0
	return cmd, nil
}
