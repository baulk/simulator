package simulator

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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
	return de
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
func (de *Derivator) InsertEnv(key, val string) {
	if strings.EqualFold(key, "PATH") {
		de.paths = append([]string{val}, de.paths...)
		return
	}
	if v, ok := de.envblocks[key]; ok {
		de.envblocks[key] = StrCat(val, string(os.PathListSeparator), v)
		return
	}
	de.envblocks[key] = val
}

// AppendEnv append
func (de *Derivator) AppendEnv(key, val string) {
	if strings.EqualFold(key, "PATH") {
		de.paths = append(de.paths, val)
		return
	}
	if v, ok := de.envblocks[key]; ok {
		de.envblocks[key] = StrCat(v, string(os.PathListSeparator), val)
		return
	}
	de.envblocks[key] = val
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
	if k == "" || v == "" {
		return errors.New("empty env k/v input")
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
func (de *Derivator) EraseEnv(k string) {
	delete(de.envblocks, k)
}

// GetEnv env
func (de *Derivator) GetEnv(k string) string {
	if v, ok := de.envblocks[k]; ok {
		return v
	}
	return os.Getenv(k)
}

// ExpandEnv env
func (de *Derivator) ExpandEnv(s string) string {
	return os.Expand(s, de.GetEnv)
}
