package pi

import (
	"errors"
	"io/fs"
	"net/http"
	"path"
	"strings"
)

func openFS(root http.FileSystem, upath string) (http.File, fs.FileInfo, error) {
	f, err := root.Open(upath)
	if err != nil {
		return nil, nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return nil, nil, err
	}
	return f, fi, nil
}

// FileServer returns a HTTP handler for serving files from within root.
// If the requested files are not exist, then send defaultsFile to client.
func FileServer(root http.FileSystem, defaultsFile string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upath := r.URL.Path
		upath = strings.TrimPrefix(upath, "/")
		upath = path.Clean(upath)

		f, fi, err := openFS(root, upath)
		if fi != nil && fi.IsDir() {
			f.Close()
			f, fi, err = openFS(root, path.Join(upath, "index.html"))
		}
		if errors.Is(err, fs.ErrNotExist) {
			f, fi, err = openFS(root, defaultsFile)
		}
		if err != nil {
			if errors.Is(err, fs.ErrPermission) {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte("403 Forbidden"))
				return
			}
			if errors.Is(err, fs.ErrNotExist) {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte("404 page not found"))
				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("500 Internal Server Error"))
			return
		}

		http.ServeContent(w, r, path.Base(upath), fi.ModTime(), f)
		f.Close()
	})
}
