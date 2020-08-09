package simulator

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// error
var (
	ErrKeyCannotEmpty = errors.New("key cannot be empty")
)

// Derivator expand env engine
type Derivator struct {
	envblocks map[string]string
	paths     []string
}

// NewDerivator create env derivative
func NewDerivator() *Derivator {
	de := &Derivator{
		envblocks: make(map[string]string),
	}
	pv := filepath.SplitList(os.Getenv("PATH"))
	de.paths = make([]string, 0, len(pv))
	for _, p := range pv {
		cp := filepath.Clean(p)
		if !de.PathListExists(cp) {
			de.paths = append(de.paths, p)
		}
	}
	for _, s := range os.Environ() {
		i := strings.IndexByte(s, '=')
		if i == -1 {
			continue
		}
		key := s[0:i]
		if strings.EqualFold(key, "PATH") {
			continue
		}
		de.envblocks[key] = s[i+1:]
	}
	if runtime.GOOS == "windows" {
		de.envblocks["home"] = os.Getenv("USERPROFILE")
	}
	return de
}

// initializeCleanupEnv todo
func (de *Derivator) initializeCleanupEnv() {
	for _, e := range allowedEnv {
		if v, b := os.LookupEnv(e); b {
			de.envblocks[e] = v
		}
	}
	if runtime.GOOS == "windows" {
		de.envblocks["home"] = os.Getenv("USERPROFILE")
	}
}

// NewCleanupDerivator todo
func NewCleanupDerivator() *Derivator {
	de := &Derivator{
		envblocks: make(map[string]string),
	}
	de.initializeCleanupEnv()
	de.paths = MakeCleanupPath()
	return de
}

// InsertEnv insert
func (de *Derivator) InsertEnv(key, val string) error {
	if key == "" {
		return ErrKeyCannotEmpty
	}
	if strings.EqualFold(key, "PATH") {
		de.paths = append([]string{val}, de.paths...)
		return nil
	}
	if v, ok := de.envblocks[key]; ok {
		de.envblocks[key] = StrCat(val, string(os.PathListSeparator), v)
		return nil
	}
	de.envblocks[key] = val
	return nil
}

// AppendEnv append
func (de *Derivator) AppendEnv(key, val string) error {
	if key == "" {
		return ErrKeyCannotEmpty
	}
	if strings.EqualFold(key, "PATH") {
		de.paths = append(de.paths, val)
		return nil
	}
	if v, ok := de.envblocks[key]; ok {
		de.envblocks[key] = StrCat(v, string(os.PathListSeparator), val)
		return nil
	}
	de.envblocks[key] = val
	return nil
}

// AddBashCompatible $0~$9
func (de *Derivator) AddBashCompatible() {
	for i := 0; i < len(os.Args); i++ {
		de.envblocks[strconv.Itoa(i)] = os.Args[i]
	}
	de.envblocks["$"] = strconv.Itoa(os.Getpid())
}

// SetEnv env
func (de *Derivator) SetEnv(k, v string) error {
	if k == "" {
		return ErrKeyCannotEmpty
	}
	de.envblocks[k] = v
	return nil
}

// Environ create new environ block
func (de *Derivator) Environ() []string {
	ev := make([]string, 0, len(de.envblocks)+1)
	for k, v := range de.envblocks {
		ev = append(ev, StrCat(k, "=", v))
	}
	ev = append(ev, StrCat("PATH=", strings.Join(de.paths, string(os.PathListSeparator))))
	return ev
}

// EraseEnv k
func (de *Derivator) EraseEnv(key string) error {
	if key == "" {
		return ErrKeyCannotEmpty
	}
	delete(de.envblocks, key)
	return nil
}

// GetEnv env
func (de *Derivator) GetEnv(key string) string {
	if key == "" {
		return ""
	}
	if v, ok := de.envblocks[key]; ok {
		return v
	}
	return ""
}

// ExpandEnv env
func (de *Derivator) ExpandEnv(s string) string {
	return os.Expand(s, de.GetEnv)
}
