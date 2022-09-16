package pi

type Result[T any] struct {
	Data         T      `json:"data"`
	Error        string `json:"error"`
	ErrorMessage string `json:"error_message"`
	Page         int    `json:"page"`
	PageSize     int    `json:"page_size"`
	Total        int    `json:"total"`
}
