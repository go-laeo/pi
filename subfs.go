package pi

import (
	"net/http"
	"os"
	"path"
)

type subfs struct {
	root     http.FileSystem
	dir      string
	defaults string
}

func (s *subfs) Open(name string) (http.File, error) {
	f, err := s.root.Open(path.Join(s.dir, name))
	if os.IsNotExist(err) && s.defaults != "" {
		f, err = s.root.Open(path.Join(s.dir, s.defaults))
	}
	return f, err
}

// Sub returns an new http.FileSystem that reads file from
// the sub directory of root.
func Sub(root http.FileSystem, dir string) http.FileSystem {
	return &subfs{root: root, dir: dir}
}

// OverrideNotFound open file from root, if the file does not
// exists, then try open defaults.
func OverrideNotFound(root http.FileSystem, defaults string) http.FileSystem {
	return &subfs{root: root, defaults: defaults}
}
