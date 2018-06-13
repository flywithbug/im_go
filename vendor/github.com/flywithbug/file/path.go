package file

import (
	"path/filepath"
	"os"
	"strings"
	"path"
)




// CurrentDir return current commond path
func CurrentDir() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	return dir, err
}

// AbsPath return p when p is a absolute path.
// If p is not absolute path and dir is a  absolute path
// then return path.Join(dir, p).
// Else get current commond path as currentPath, reutrn path.Join(currentPath, dir, p)
func AbsPath(dir string, p string) string {
	p = strings.TrimSpace(p)
	if path.IsAbs(p) {
		return p
	}

	dir = strings.TrimSpace(dir)
	if dir != "" && path.IsAbs(dir) {
		return path.Join(dir, p)
	}

	if d, err := CurrentDir(); err == nil {
		return path.Join(d, dir, p)
	}

	return p
}

// EnumeratePath enumerate the path
func EnumeratePath(p string, f func(surplus string, current string, stop *bool)) {
	if f == nil {
		return
	}
	stop := false
	for p != "" {
		c := path.Base(p)
		p = path.Dir(p)
		f(p, c, &stop)
		if stop {
			break
		}
	}
}