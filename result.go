package pi

type Result[T any] struct {
	Data T    `json:"data"`
	OK   bool `json:"ok"`
}

type ErrorResult struct {
	Error        string `json:"error"`
	ErrorMessage string `json:"error_message"`
	OK           bool   `json:"ok"`
}

type LengthResult[T any] struct {
	Data     []T  `json:"data"`
	Page     int  `json:"page"`
	PageSize int  `json:"page_size"`
	Total    int  `json:"total"`
	OK       bool `json:"ok"`
}
