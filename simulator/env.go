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
	home   string
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
		if strings.EqualFold(key, "HOME") || strings.EqualFold(key, "USERPROFILE") {
			sm.home = s[i+1:]
			continue
		}
		if strings.EqualFold(key, "PATH") {
			continue
		}
		sm.envmap[key] = s[i+1:]
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
		sm.home = os.Getenv("USERPROFILE")
	} else {
		sm.home = os.Getenv("HOME")
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
	ev := make([]string, 0, len(sm.envmap)+3)
	var hasUserProfile bool
	for k, v := range sm.envmap {
		if strings.EqualFold(k, "USERPROFILE") {
			hasUserProfile = true
		}
		ev = append(ev, StrCat(k, "=", v))
	}
	ev = append(ev, StrCat("HOME=", sm.home))
	if runtime.GOOS == "windows" && !hasUserProfile {
		ev = append(ev, StrCat("USERPROFILE=", sm.home))
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

// LookupEnv lookup env
func (sm *Simulator) LookupEnv(key string) (string, bool) {
	if key == "" {
		return "", false
	}
	if v, ok := sm.envmap[key]; ok {
		return v, true
	}
	return "", false
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

// SetHome simulate home
func (sm *Simulator) SetHome(home string) string {
	h := sm.home
	sm.home = home
	return h
}

// ExpandEnv env
func (sm *Simulator) ExpandEnv(s string) string {
	return os.Expand(s, sm.GetEnv)
}

// Paths return paths
func (sm *Simulator) Paths() []string {
	return sm.paths
}

// Home return home dir
func (sm *Simulator) Home() string {
	return sm.home
}
