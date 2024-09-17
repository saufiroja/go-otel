package responses

type Response[T any] struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    T      `json:"data"`
}

func NewResponse[T any](message string, code int, data T) *Response[T] {
	return &Response[T]{
		Message: message,
		Code:    code,
		Data:    data,
	}
}
