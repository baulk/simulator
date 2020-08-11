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

// Simulator expand env engine
type Simulator struct {
	envmap map[string]string
	paths  []string
}

// NewSimulator create env simulator
func NewSimulator() *Simulator {
	sm := &Simulator{
		envmap: make(map[string]string),
	}
	pv := filepath.SplitList(os.Getenv("PATH"))
	sm.paths = make([]string, 0, len(pv))
	for _, p := range pv {
		cp := filepath.Clean(p)
		if !sm.PathListExists(cp) {
			sm.paths = append(sm.paths, p)
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
		sm.envmap[key] = s[i+1:]
	}
	if runtime.GOOS == "windows" {
		sm.envmap["home"] = os.Getenv("USERPROFILE")
	}
	return sm
}

// initializeCleanupEnv todo
func (sm *Simulator) initializeCleanupEnv() {
	for _, e := range allowedEnv {
		if v, b := os.LookupEnv(e); b {
			sm.envmap[e] = v
		}
	}
	if runtime.GOOS == "windows" {
		sm.envmap["home"] = os.Getenv("USERPROFILE")
	}
}

// NewCleanupSimulator todo
func NewCleanupSimulator() *Simulator {
	sm := &Simulator{
		envmap: make(map[string]string),
	}
	sm.initializeCleanupEnv()
	sm.paths = MakeCleanupPath()
	return sm
}

// InsertEnv insert
func (sm *Simulator) InsertEnv(key, val string) error {
	if key == "" {
		return ErrKeyCannotEmpty
	}
	if strings.EqualFold(key, "PATH") {
		sm.paths = append([]string{val}, sm.paths...)
		return nil
	}
	if v, ok := sm.envmap[key]; ok {
		sm.envmap[key] = StrCat(val, string(os.PathListSeparator), v)
		return nil
	}
	sm.envmap[key] = val
	return nil
}

// AppendEnv append
func (sm *Simulator) AppendEnv(key, val string) error {
	if key == "" {
		return ErrKeyCannotEmpty
	}
	if strings.EqualFold(key, "PATH") {
		sm.paths = append(sm.paths, val)
		return nil
	}
	if v, ok := sm.envmap[key]; ok {
		sm.envmap[key] = StrCat(v, string(os.PathListSeparator), val)
		return nil
	}
	sm.envmap[key] = val
	return nil
}

// AddBashCompatible $0~$9
func (sm *Simulator) AddBashCompatible() {
	for i := 0; i < len(os.Args); i++ {
		sm.envmap[strconv.Itoa(i)] = os.Args[i]
	}
	sm.envmap["$"] = strconv.Itoa(os.Getpid())
}

// SetEnv env
func (sm *Simulator) SetEnv(k, v string) error {
	if k == "" {
		return ErrKeyCannotEmpty
	}
	sm.envmap[k] = v
	return nil
}

// Environ create new environ block
func (sm *Simulator) Environ() []string {
	ev := make([]string, 0, len(sm.envmap)+1)
	for k, v := range sm.envmap {
		ev = append(ev, StrCat(k, "=", v))
	}
	ev = append(ev, StrCat("PATH=", strings.Join(sm.paths, string(os.PathListSeparator))))
	return ev
}

// EraseEnv k
func (sm *Simulator) EraseEnv(key string) error {
	if key == "" {
		return ErrKeyCannotEmpty
	}
	delete(sm.envmap, key)
	return nil
}

// GetEnv env
func (sm *Simulator) GetEnv(key string) string {
	if key == "" {
		return ""
	}
	if v, ok := sm.envmap[key]; ok {
		return v
	}
	return ""
}

// ExpandEnv env
func (sm *Simulator) ExpandEnv(s string) string {
	return os.Expand(s, sm.GetEnv)
}
