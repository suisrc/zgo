package helper

// Page 分页数据
type Page struct {
	PageNo   int         `json:"pageNo,omitempty"`   // 页索引
	PageSize int         `json:"pageSize,omitempty"` // 页条数
	Total    int         `json:"total,omitempty"`    // 总条数
	List     interface{} `json:"list,omitempty"`     // 数据
}
