package req

// @Description 分页查询结构
type PageReq struct {
	// 页码
	Page int `json:"page" form:"page"`
	// 每页大小
	PageSize int `json:"page_size" form:"page_size"`
}

func (r *PageReq) Offset() int {
	return (r.Page - 1) * r.PageSize
}

// @Description 分页查询请求
type PageUserReq struct {
	PageReq
	// 用户名
	UserName string `json:"user_name" form:"user_name"`
}
