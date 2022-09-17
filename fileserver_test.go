package pi

import (
	"bytes"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"testing"
	"testing/fstest"
	"time"
)

func TestFileServer(t *testing.T) {
	root := fstest.MapFS{
		"web/dist/index.html":  &fstest.MapFile{Data: []byte("index.html"), Mode: fs.ModePerm, ModTime: time.Now()},
		"web/dist/css/app.css": &fstest.MapFile{Data: []byte("app.css"), Mode: fs.ModePerm, ModTime: time.Now()},
	}

	h := FileServer(http.FS(root), "web/dist/index.html")

	tests := []struct {
		name      string
		target    string
		wantBytes []byte
		wantCode  int
	}{
		{
			name:      "request / should ok",
			target:    "/",
			wantCode:  200,
			wantBytes: []byte("index.html"),
		},
		{
			name:      "request /web/dist/index.html should ok",
			target:    "/web/dist/index.html",
			wantCode:  200,
			wantBytes: []byte("index.html"),
		},
		{
			name:      "request /index.html should got defaultsFile",
			target:    "/index.html",
			wantCode:  200,
			wantBytes: []byte("index.html"),
		},
		{
			name:      "request /web/dist/css/app.css should ok",
			target:    "/web/dist/css/app.css",
			wantCode:  200,
			wantBytes: []byte("app.css"),
		},
		{
			name:      "request /web/dist/css should got defaultsFile",
			target:    "/web/dist/css",
			wantCode:  200,
			wantBytes: []byte("index.html"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, tt.target, nil)
			h(createContext(w, r, nil, nil))
			if w.Code != tt.wantCode {
				t.Fatalf("%s want code = %d, got %d", tt.name, tt.wantCode, w.Code)
			}
			if !bytes.Equal(w.Body.Bytes(), tt.wantBytes) {
				t.Fatalf("%s want body = %s, got %s", tt.name, tt.wantBytes, w.Body.Bytes())
			}
		})
	}
}

func BenchmarkFileServer_ServeHTTP(b *testing.B) {
	root := fstest.MapFS{
		"web/dist/index.html":  &fstest.MapFile{Data: []byte("index.html"), Mode: fs.ModePerm, ModTime: time.Now()},
		"web/dist/css/app.css": &fstest.MapFile{Data: []byte("app.css"), Mode: fs.ModePerm, ModTime: time.Now()},
	}

	h := FileServer(http.FS(root), "web/dist/index.html")
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := createContext(w, r, nil, nil)
		h(ctx)
	}
}
