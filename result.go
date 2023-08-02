package pi

type Result[T any] struct {
	Data T `json:"data"`
}

type ErrorResult struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"error_message"`
}

type LengthResult[T any] struct {
	Data     []T `json:"data"`
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
}
