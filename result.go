package pi

type Result[T any] struct {
	Data         T      `json:"data,omitempty"`
	Error        string `json:"error,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	OK           bool   `json:"ok"`
}

type LengthResult[T any] struct {
	Data         T      `json:"data,omitempty"`
	Error        string `json:"error,omitempty"`
	ErrorMessage string `json:"error_message,omitempty"`
	Page         int    `json:"page,omitempty"`
	PageSize     int    `json:"page_size,omitempty"`
	Total        int    `json:"total,omitempty"`
	OK           bool   `json:"ok"`
}
