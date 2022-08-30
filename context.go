package ezy

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"time"
)

type Context interface {
	context.Context
	http.ResponseWriter

	WithContext(ctx context.Context) Context

	Raw() (w http.ResponseWriter, r *http.Request)
	Query(field string, defaults ...string) string
	Form(field string, defaults ...string) string
	File(field string) (multipart.File, *multipart.FileHeader, error)
	FileSet(field string) []*multipart.FileHeader
	Cookie(name string) (*http.Cookie, error)
	HeaderMap() http.Header
	HeaderValue(name string) string

	Format(v any) error
	Json(v any) error
	Redirect(to string, code ...int) error
}

type ctx struct {
	w http.ResponseWriter
	r *http.Request
	b *bytes.Buffer
	c int
}

var _ Context = (*ctx)(nil)

// Deadline returns the time when work done on behalf of this context
// should be canceled. Deadline returns ok==false when no deadline is
// set. Successive calls to Deadline return the same results.
func (c *ctx) Deadline() (deadline time.Time, ok bool) {
	return c.r.Context().Deadline()
}

// Done returns a channel that's closed when work done on behalf of this
// context should be canceled. Done may return nil if this context can
// never be canceled. Successive calls to Done return the same value.
// The close of the Done channel may happen asynchronously,
// after the cancel function returns.
//
// WithCancel arranges for Done to be closed when cancel is called;
// WithDeadline arranges for Done to be closed when the deadline
// expires; WithTimeout arranges for Done to be closed when the timeout
// elapses.
//
// Done is provided for use in select statements:
//
//	// Stream generates values with DoSomething and sends them to out
//	// until DoSomething returns an error or ctx.Done is closed.
//	func Stream(ctx context.Context, out chan<- Value) error {
//		for {
//			v, err := DoSomething(ctx)
//			if err != nil {
//				return err
//			}
//			select {
//			case <-ctx.Done():
//				return ctx.Err()
//			case out <- v:
//			}
//		}
//	}
//
// See https://blog.golang.org/pipelines for more examples of how to use
// a Done channel for cancellation.
func (c *ctx) Done() <-chan struct{} {
	return c.r.Context().Done()
}

// If Done is not yet closed, Err returns nil.
// If Done is closed, Err returns a non-nil error explaining why:
// Canceled if the context was canceled
// or DeadlineExceeded if the context's deadline passed.
// After Err returns a non-nil error, successive calls to Err return the same error.
func (c *ctx) Err() error {
	return c.r.Context().Err()
}

// Value returns the value associated with this context for key, or nil
// if no value is associated with key. Successive calls to Value with
// the same key returns the same result.
//
// Use context values only for request-scoped data that transits
// processes and API boundaries, not for passing optional parameters to
// functions.
//
// A key identifies a specific value in a Context. Functions that wish
// to store values in Context typically allocate a key in a global
// variable then use that key as the argument to context.WithValue and
// Context.Value. A key can be any type that supports equality;
// packages should define keys as an unexported type to avoid
// collisions.
//
// Packages that define a Context key should provide type-safe accessors
// for the values stored using that key:
//
//	// Package user defines a User type that's stored in Contexts.
//	package user
//
//	import "context"
//
//	// User is the type of value stored in the Contexts.
//	type User struct {...}
//
//	// key is an unexported type for keys defined in this package.
//	// This prevents collisions with keys defined in other packages.
//	type key int
//
//	// userKey is the key for user.User values in Contexts. It is
//	// unexported; clients use user.NewContext and user.FromContext
//	// instead of using this key directly.
//	var userKey key
//
//	// NewContext returns a new Context that carries value u.
//	func NewContext(ctx context.Context, u *User) context.Context {
//		return context.WithValue(ctx, userKey, u)
//	}
//
//	// FromContext returns the User value stored in ctx, if any.
//	func FromContext(ctx context.Context) (*User, bool) {
//		u, ok := ctx.Value(userKey).(*User)
//		return u, ok
//	}
func (c *ctx) Value(key any) any {
	return c.r.Context().Value(key)
}

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
func (c *ctx) Header() http.Header {
	return c.w.Header()
}

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
func (c *ctx) Write(b []byte) (int, error) {
	return c.b.Write(b)
}

// WriteHeader sends an HTTP response header with the provided
// status code.
//
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes or 1xx informational responses.
//
// The provided code must be a valid HTTP 1xx-5xx status code.
// Any number of 1xx headers may be written, followed by at most
// one 2xx-5xx header. 1xx headers are sent immediately, but 2xx-5xx
// headers may be buffered. Use the Flusher interface to send
// buffered data. The header map is cleared when 2xx-5xx headers are
// sent, but not with 1xx headers.
//
// The server will automatically send a 100 (Continue) header
// on the first read from the request body if the request has
// an "Expect: 100-continue" header.
func (c *ctx) WriteHeader(statusCode int) {
	c.c = statusCode
}

func (c *ctx) WithContext(ctx context.Context) Context {
	c.r = c.r.WithContext(ctx)
	return c
}

func (c *ctx) Raw() (w http.ResponseWriter, r *http.Request) {
	return c, c.r
}

func (c *ctx) Query(field string, defaults ...string) string {
	if c.r.URL.Query().Has(field) {
		return c.r.URL.Query().Get(field)
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

func (c *ctx) Form(field string, defaults ...string) string {
	v := c.r.FormValue(field)
	if v != "" {
		return v
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

func (c *ctx) File(field string) (multipart.File, *multipart.FileHeader, error) {
	return c.r.FormFile(field)
}

func (c *ctx) FileSet(field string) []*multipart.FileHeader {
	c.r.ParseMultipartForm(1024 * 1024 * 1024)
	return c.r.MultipartForm.File[field]
}

func (c *ctx) Cookie(name string) (*http.Cookie, error) {
	return c.r.Cookie(name)
}

func (c *ctx) HeaderMap() http.Header {
	return c.r.Header
}

func (c *ctx) HeaderValue(name string) string {
	return c.r.Header.Get(name)
}

func (c *ctx) Format(v any) error {
	panic("not implemented") // TODO: Implement
}

func (c *ctx) Json(v any) error {
	return json.NewEncoder(c).Encode(v)
}

func (c *ctx) Redirect(to string, code ...int) error {
	if len(code) == 0 {
		http.Redirect(c.w, c.r, to, http.StatusTemporaryRedirect)
		return nil
	}

	http.Redirect(c.w, c.r, to, code[0])
	return nil
}
