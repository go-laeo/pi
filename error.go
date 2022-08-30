package ezy

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

func NewError(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func (e *Error) Error() string {
	return e.Message
}
