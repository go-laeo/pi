package pi

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"net/http"
	"testing"
	"testing/fstest"
	"time"
)

func TestSub(t *testing.T) {
	root := fstest.MapFS{
		"web/dist/index.html":  &fstest.MapFile{Data: []byte("index.html"), Mode: fs.ModePerm, ModTime: time.Now()},
		"web/dist/css/app.css": &fstest.MapFile{Data: []byte("app.css"), Mode: fs.ModePerm, ModTime: time.Now()},
	}

	tests := []struct {
		name            string
		dir             string
		open            string
		wantBytes       []byte
		wantDir         bool
		wantErrNotExist bool
	}{
		{
			name:      "open web/dist/index.html by /index.html should ok",
			dir:       "web/dist",
			open:      "/index.html",
			wantBytes: []byte("index.html"),
		},
		{
			name:      "open web/dist/index.html by index.html should ok",
			dir:       "web/dist",
			open:      "index.html",
			wantBytes: []byte("index.html"),
		},
		{
			name:    "open / should got a dir",
			dir:     "web/dist",
			open:    "/",
			wantDir: true,
		},
		{
			name:            "open unknown.txt should got ErrNotFound",
			dir:             "web/dist",
			open:            "unknown.txt",
			wantErrNotExist: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := Sub(http.FS(root), tt.dir)
			f, err := h.Open(tt.open)
			if err != nil {
				if !errors.Is(err, fs.ErrNotExist) || !tt.wantErrNotExist {
					t.Fatalf("%s got error = %v", tt.name, err)
				}
			}
			if tt.wantDir {
				fi, err := f.Stat()
				if err != nil {
					t.Fatalf("%s got error = %v", tt.name, err)
				}
				if !fi.IsDir() {
					t.Fatalf("%s got a file", tt.name)
				}
			}
			if tt.wantBytes != nil {
				b, err := io.ReadAll(f)
				if err != nil {
					t.Fatalf("%s got error = %v", tt.name, err)
				}

				f.Close()

				if !bytes.Equal(b, tt.wantBytes) {
					t.Fatalf("%s want bytes = %s, got = %s", tt.name, tt.wantBytes, b)
				}
			}
		})
	}
}
