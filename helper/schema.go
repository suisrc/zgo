package helper

const (
	// ShowNone 静音
	ShowNone = 0
	// ShowWarn 消息警告
	ShowWarn = 1
	// ShowError 消息错误
	ShowError = 2
	// ShowNotify 通知；
	ShowNotify = 4
	// ShowPage 页
	ShowPage = 9
)

// Success 正常请求结构体
type Success struct {
	Success bool        `json:"success"`        // 请求成功, false
	Data    interface{} `json:"data,omitempty"` // 响应数据
	TraceID string      `json:"traceId"`        // 方便进行后端故障排除：唯一的请求ID
}

func (e *Success) Error() string {
	return "success"
}

// ErrorInfo 异常的请求结果体
type ErrorInfo struct {
	Success      bool        `json:"success"`        // 请求成功, false
	Data         interface{} `json:"data,omitempty"` // 响应数据
	ErrorCode    string      `json:"errorCode"`      // 错误代码
	ErrorMessage string      `json:"errorMessage"`   // 向用户显示消息
	ShowType     int         `json:"showType"`       //错误显示类型：0静音； 1条消息警告； 2消息错误； 4通知； 9页
	TraceID      string      `json:"traceId"`        // 方便进行后端故障排除：唯一的请求ID
	//Status       int         `json:"-"`
}

func (e *ErrorInfo) Error() string {
	return e.ErrorMessage
}

// ErrorRedirect 重定向
type ErrorRedirect struct {
	Status   int    // http.StatusSeeOther
	State    string // 状态, 用户还原缓存现场
	Location string
}

func (e *ErrorRedirect) Error() string {
	return "Redirect: " + e.Location
}

// ErrorNone 返回值已经被处理,无返回值
type ErrorNone struct {
}

func (e *ErrorNone) Error() string {
	return "http none"
}

// PaginationResult 响应列表数据
//type PaginationResult struct {
//	list  interface{} `json:"list"`
//	total int         `json:"total"`
//	sign  string      `json:"sign,omitempty" binding:"required"`
//}

// PaginationParam 分页查询条件
type PaginationParam struct {
	PageSign  string `query:"pageSign"`                              // 请求参数, total | list | both
	PageNo    uint   `query:"pageNo,default=1"`                      // 当前页
	PageSize  uint   `query:"pageSize,default=20" binding:"max=100"` // 页大小
	PageTotal uint   `query:"pageTotal"`                             // 上次统计的数据条数
}
