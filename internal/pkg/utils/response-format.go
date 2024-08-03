package utils

type Response[T any] struct {
	Data     T           `json:"data"`
	Code     string      `json:"code"`
	Message  string      `json:"message"`
	Status   int         `json:"status"`
	Metadata interface{} `json:"meta"`
}

func NewResponse[T any](data T, message string, status int, meta interface{}) *Response[T] {
	return &Response[T]{
		Data:     data,
		Message:  message,
		Status:   status,
		Metadata: meta,
	}
}
