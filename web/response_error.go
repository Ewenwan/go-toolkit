package web

// ErrorResponse Error Response
type ErrorResponse struct {
	response *Response
	original string
	code     int
}

// NewErrorResponse Create error response
func NewErrorResponse(response *Response, res string, code int) *ErrorResponse {
	return &ErrorResponse{
		response: response,
		original: res,
		code:     code,
	}
}

// WithCode set response code and return itself
func (resp *ErrorResponse) WithCode(code int) *ErrorResponse {
	resp.code = code
	return resp
}

// CreateResponse 创建响应内容
func (resp *ErrorResponse) CreateResponse() error {
	resp.response.SendError(resp.code, resp.original)
	resp.response.Flush()
	return nil
}
