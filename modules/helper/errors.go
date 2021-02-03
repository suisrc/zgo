package helper

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	gi18n "github.com/suisrc/gin-i18n"
	"github.com/suisrc/zgo/modules/logger"

	"github.com/gin-gonic/gin"
)

// H h -> map
type H map[string]interface{}

// ErrorModel 异常模型
type ErrorModel struct {
	Status       int
	ShowType     int
	ErrorMessage *i18n.Message
	ErrorArgs    map[string]interface{}
}

func (a *ErrorModel) Error() string {
	return fmt.Sprintf("[%d]%s:%s", a.Status, a.ErrorMessage.ID, a.ErrorMessage.Other)
}

// 定义错误
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Status/405
var (
	Err400BadRequest       = &ErrorModel{Status: 400, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: "ERR-BAD-REQUEST", Other: "请求发生错误"}}
	Err401Unauthorized     = &ErrorModel{Status: 401, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: "ERR-UNAUTHORIZED", Other: "用户没有权限（令牌、用户名、密码错误）"}}
	Err403Forbidden        = &ErrorModel{Status: 403, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: "ERR-FORBIDDEN", Other: "用户得到授权，但是访问是被禁止的"}}
	Err404NotFound         = &ErrorModel{Status: 404, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: "ERR-NOT-FOUND", Other: "发出的请求针对的是不存在的记录，服务器没有进行操作"}}
	Err405MethodNotAllowed = &ErrorModel{Status: 405, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: "ERR-METHOD-NOT-ALLOWED", Other: "请求的方法不允许"}}
	Err406NotAcceptable    = &ErrorModel{Status: 406, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: "ERR-NOT-ACCEPTABLE", Other: "请求的格式不可得"}}
	Err429TooManyRequests  = &ErrorModel{Status: 429, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: "ERR-TOO-MANY-REQUESTS", Other: "请求次数过多"}}
	Err456TokenExpired     = &ErrorModel{Status: 456, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: "ERR-TOKEN-EXPIRED", Other: "请求令牌已过期"}}
	Err500InternalServer   = &ErrorModel{Status: 500, ShowType: ShowWarn, ErrorMessage: &i18n.Message{ID: "ERR-INTERNAL-SERVER", Other: "服务器发生错误"}}
)

// NewError 包装响应错误
func NewError(ctx *gin.Context, showType int, emsg *i18n.Message, args map[string]interface{}) *ErrorInfo {
	res := &ErrorInfo{
		Success:      false,
		ErrorCode:    emsg.ID,
		ErrorMessage: gi18n.FormatMessage(ctx, emsg, args),
		ShowType:     showType,
		TraceID:      GetTraceID(ctx),
		//Status:       http.StatusOK,
	}
	return res
}

// New0Error 包装响应错误, 没有参数
func New0Error(ctx *gin.Context, showType int, emsg *i18n.Message) *ErrorInfo {
	return NewError(ctx, showType, emsg, nil)
}

// NewWrapError 包装响应错误
func NewWrapError(ctx *gin.Context, em *ErrorModel) *ErrorInfo {
	res := &ErrorInfo{
		Success:      false,
		ErrorCode:    em.ErrorMessage.ID,
		ErrorMessage: gi18n.FormatMessage(ctx, em.ErrorMessage, em.ErrorArgs),
		ShowType:     em.ShowType,
		TraceID:      GetTraceID(ctx),
		//Status:       em.Status,
	}
	return res
}

// NewSuccess 包装响应结果
func NewSuccess(ctx *gin.Context, data interface{}) *Success {
	res := &Success{
		Success: true,
		Data:    data,
		TraceID: GetTraceID(ctx),
	}
	return res
}

// Wrap400Response 无法解析异常
func Wrap400Response(ctx *gin.Context, err error) *ErrorModel {
	return &ErrorModel{
		Status:       400,
		ShowType:     ShowWarn,
		ErrorMessage: &i18n.Message{ID: "ERR-BAD-REQUEST-X", Other: "解析请求参数发生错误 - {{.error}}"},
		ErrorArgs:    map[string]interface{}{"error": logger.ErrorWW(err)},
	}
}
