// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package simulator

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func findExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return os.ErrPermission
}

// LookPath look path
func (sm *Simulator) LookPath(file string) (string, error) {
	// NOTE(rsc): I wish we could use the Plan 9 behavior here
	// (only bypass the path if file begins with / or ./ or ../)
	// but that would not match all the Unix shells.

	if strings.Contains(file, "/") {
		err := findExecutable(file)
		if err == nil {
			return file, nil
		}
		return "", &Error{file, err}
	}
	for _, dir := range sm.paths {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		path := filepath.Join(dir, file)
		if err := findExecutable(path); err == nil {
			return path, nil
		}
	}
	return "", &Error{file, exec.ErrNotFound}
}

// PathListExists path exists
func (sm *Simulator) PathListExists(p string) bool {
	for _, s := range sm.paths {
		if s == p {
			return true
		}
	}
	return false
}

// MakeCleanupPath make cleanup
func MakeCleanupPath() []string {
	return []string{
		"/usr/local/sbin",
		"/usr/local/bin",
		"/usr/sbin",
		"/usr/bin",
		"/sbin",
		"/bin",
	}
}

var allowedEnv = []string{
	"HOSTTYPE",
	"LANG",
	"TERM",
	"NAME",
	"USER",
	"LONGNAME",
	"SHELL",
	"TZ",
	"LD_LIBRARY_PATH",
	// Enables proxy information to be passed to Curl, the unsmrlying download
	// library in cmake.exe
	"http_proxy",
	"https_proxy",
	// Environment variables to tell git to use custom SSH executable or command
	"GIT_SSH",
	"GIT_SSH_COMMAND",
	// Environment variables neesmd for ssh-agent based authentication
	"SSH_AUTH_SOCK",
	"SSH_AGENT_PID",
}
