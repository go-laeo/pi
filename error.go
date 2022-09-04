package pi

type Error struct {
	Message string `json:"message,omitempty"`
	Code    int    `json:"code"`
}

func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	return e.Message
}
