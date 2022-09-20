package pi

import (
	"context"
	"encoding/json"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
)

type Context interface {
	Context() context.Context
	SetContext(ctx context.Context)

	Raw() (w http.ResponseWriter, r *http.Request)
	Query(field string, defaults ...string) string
	Form(field string, defaults ...string) string

	// File returns first uploaded file by field.
	File(field string) (multipart.File, *multipart.FileHeader, error)

	// FileSet gets all uploaded files by field from underlying request,
	FileSet(field string) []*multipart.FileHeader

	// Cookie returns the named cookie provided in the request or
	// ErrNoCookie if not found.
	// If multiple cookies match the given name, only one cookie will
	// be returned.
	Cookie(name string) (*http.Cookie, error)

	// Get gets the first value associated with the given key from request header. If
	// there are no values associated with the key, Get returns "".
	// It is case insensitive; textproto.CanonicalMIMEHeaderKey is
	// used to canonicalize the provided key. Get assumes that all
	// keys are stored in canonical form. To use non-canonical keys,
	// access the map directly.
	Get(name string) string

	// Domain gets domain name of from request's Host field, eg. www.google.com.
	Domain() string

	URL() *url.URL

	// Param gets named route param by name, returns empty string if it does not exists.
	Param(name string) string

	// ParamValues returns all path params.
	ParamValues() url.Values

	// IP gets first client IP.
	IP() string

	// IPSet gets all client IPs if available.
	IPSet() []string

	Method() string
	Is(method string) bool

	// Header returns the header map that will be sent by
	// WriteHeader. The Header map also is the mechanism with which
	// Handlers can set HTTP trailers.
	//
	// Changing the header map after a call to WriteHeader (or
	// Write) has no effect unless the HTTP status code was of the
	// 1xx class or the modified headers are trailers.
	//
	// There are two ways to set Trailers. The preferred way is to
	// predeclare in the headers which trailers you will later
	// send by setting the "Trailer" header to the names of the
	// trailer keys which will come later. In this case, those
	// keys of the Header map are treated as if they were
	// trailers. See the example. The second way, for trailer
	// keys not known to the Handler until after the first Write,
	// is to prefix the Header map keys with the TrailerPrefix
	// constant value. See TrailerPrefix.
	//
	// To suppress automatic response headers (such as "Date"), set
	// their value to nil.
	Header() http.Header

	// Write writes the data to the connection as part of an HTTP reply.
	//
	// If WriteHeader has not yet been called, Write calls
	// WriteHeader(http.StatusOK) before writing the data. If the Header
	// does not contain a Content-Type line, Write adds a Content-Type set
	// to the result of passing the initial 512 bytes of written data to
	// DetectContentType. Additionally, if the total size of all written
	// data is under a few KB and there are no Flush calls, the
	// Content-Length header is added automatically.
	//
	// Depending on the HTTP protocol version and the client, calling
	// Write or WriteHeader may prevent future reads on the
	// Request.Body. For HTTP/1.x requests, handlers should read any
	// needed request body data before writing the response. Once the
	// headers have been flushed (due to either an explicit Flusher.Flush
	// call or writing enough data to trigger a flush), the request body
	// may be unavailable. For HTTP/2 requests, the Go HTTP server permits
	// handlers to continue to read the request body while concurrently
	// writing the response. However, such behavior may not be supported
	// by all HTTP/2 clients. Handlers should read before writing if
	// possible to maximize compatibility.
	Write([]byte) (int, error)

	// Code only responds to client with HTTP status code.
	Code(code int) error

	// Json encode v into JSON string then writes to client.
	Json(v any) error

	// Text writes v to client as plain text.
	Text(v string) error
	Redirect(to string, code ...int) error
	SetCookie(c *http.Cookie)
}

var _ Context = (*_ctx)(nil)

type _ctx struct {
	w http.ResponseWriter
	r *http.Request
	p url.Values
}

func createContext(w http.ResponseWriter, r *http.Request, cap url.Values) Context {
	return &_ctx{
		w: w,
		r: r,
		p: cap,
	}
}

func (c *_ctx) Header() http.Header {
	return c.w.Header()
}

func (c *_ctx) Write(b []byte) (int, error) {
	return c.w.Write(b)
}

func (c *_ctx) Context() context.Context {
	return c.r.Context()
}

func (c *_ctx) SetContext(ctx context.Context) {
	c.r = c.r.WithContext(ctx)
}

func (c *_ctx) Raw() (w http.ResponseWriter, r *http.Request) {
	return c.w, c.r
}

func (c *_ctx) Query(field string, defaults ...string) string {
	if c.r.URL.Query().Has(field) {
		return c.r.URL.Query().Get(field)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

func (c *_ctx) Form(field string, defaults ...string) string {
	v := c.r.FormValue(field)
	if v != "" {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

func (c *_ctx) File(field string) (multipart.File, *multipart.FileHeader, error) {
	return c.r.FormFile(field)
}

func (c *_ctx) FileSet(field string) []*multipart.FileHeader {
	c.r.ParseMultipartForm(1024 * 1024 * 1024)
	return c.r.MultipartForm.File[field]
}

func (c *_ctx) Cookie(name string) (*http.Cookie, error) {
	return c.r.Cookie(name)
}

func (c *_ctx) Get(name string) string {
	return c.r.Header.Get(name)
}

func (c *_ctx) Domain() string {
	return c.r.Host
}

func (c *_ctx) URL() *url.URL {
	return c.r.URL
}

func (c *_ctx) Param(name string) string {
	return c.p.Get(name)
}

func (c *_ctx) ParamValues() url.Values {
	return c.p
}

func (c *_ctx) IP() string {
	host, _, err := net.SplitHostPort(c.r.RemoteAddr)
	if err != nil {
		return ""
	}
	return host
}

func (c *_ctx) IPSet() []string {
	host, _, err := net.SplitHostPort(c.r.RemoteAddr)
	if err != nil {
		return nil
	}
	return []string{host}
}

func (c *_ctx) Method() string {
	return c.r.Method
}

func (c *_ctx) Is(method string) bool {
	return c.Method() == method
}

// Code sends an HTTP response header with the provided
// status code.
func (c *_ctx) Code(code int) (err error) {
	c.w.WriteHeader(code)
	return
}

func (c *_ctx) Json(v any) error {
	c.Header().Set("content-type", "application/json")
	return json.NewEncoder(c).Encode(v)
}

func (c *_ctx) Text(v string) error {
	c.Header().Set("content-type", "text/plain")
	_, err := c.Write([]byte(v))
	return err
}

func (c *_ctx) Redirect(to string, code ...int) error {
	if len(code) == 0 {
		http.Redirect(c.w, c.r, to, http.StatusTemporaryRedirect)
		return nil
	}

	http.Redirect(c.w, c.r, to, code[0])
	return nil
}

func (c *_ctx) SetCookie(co *http.Cookie) {
	http.SetCookie(c.w, co)
}
