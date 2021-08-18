package response

const (
	ResponseResultSuccess = "success"
	ResponseResultError   = "error"
	ResponseResultWarning = "warning"
)

type ResponseResult struct {
	Code    int         `json:"code"`
	Type    string      `json:"type"`
	Message string      `json:"message,omitempty"`
	Url     string      `json:"url,omitempty"`
	Result  interface{} `json:"result,omitempty"`
}

type ResponsePaginationResult struct {
	Total int64       `json:"total"`
	Items interface{} `json:"items"`
}
