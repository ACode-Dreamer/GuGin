package data

// @Description 基础序列化响应
type Response struct {
	// 业务处理状态
	Success bool `json:"success"`
	// 错误编码
	ErrCode int `json:"err_code,omitempty"`
	// 数据
	Data interface{} `json:"data,omitempty"`
	// 信息
	Message string `json:"message,omitempty"`
}

// @Description 分页结构体
type Pagination struct {
	// 总条数
	Total int64 `json:"total"`
	// 数据
	Items interface{} `json:"items,omitempty"`
}

// 三位数错误编码为复用http原本含义
// 五位数错误编码为应用自定义错误
// 五开头的五位数错误编码为服务器端错误，比如数据库操作失败
// 四开头的五位数错误编码为客户端错误，有时候是客户端代码写错了，有时候是用户操作错误
const (
	// CodeCheckLogin 未登录
	CodeCheckLogin = 401
	// CodeNoRightErr 未授权访问
	CodeNoRightErr = 403
	// CodeDBError 数据库操作失败
	CodeDBError = 50001
	// CodeEncryptError 加密失败
	CodeEncryptError = 50002
	// CodeParamErr 各种奇奇怪怪的参数错误
	CodeParamErr = 40001
)

// NewErrorResponse 通用错误处理
func NewErrorResponse(errCode int, msg string) *Response {
	res := &Response{
		Success: false,
		ErrCode: errCode,
		Message: msg,
	}
	return res
}

// NewSuccessResponse 通用信息处理
func NewSuccessResponse() *Response {
	res := &Response{
		Success: true,
		Message: "操作成功",
	}
	return res
}

// NewDataResponse 通用参数处理
func NewDataResponse(data interface{}) *Response {
	res := &Response{
		Success: true,
		Data:    data,
	}
	return res
}

// NewPageResponse 通用分页参数处理
func NewPageResponse(total int64, array interface{}) *Response {
	res := &Response{
		Success: true,
		Data: &Pagination{
			Total: total,
			Items: array,
		},
	}
	return res
}

// ParamErr 各种参数错误
func ParamErr(msg string) *Response {
	if msg == "" {
		msg = "参数错误"
	}
	return NewErrorResponse(CodeParamErr, msg)
}
