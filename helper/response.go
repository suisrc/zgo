package helper

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suisrc/zgo/logger"
)

// ResSuccess 包装响应错误
// 禁止service层调用,请使用NewSuccess替换
func ResSuccess(ctx *gin.Context, v interface{}) error {
	res := NewSuccess(ctx, v)
	//ctx.JSON(http.StatusOK, res)
	//ctx.Abort()
	ResJSON(ctx, http.StatusOK, res)
	return res
}

// ResError 包装响应错误
// 禁止service层调用,请使用NewWarpError替换
func ResError(ctx *gin.Context, em *ErrorModel) error {
	res := NewWrapError(ctx, em)
	//ctx.JSON(http.StatusOK, res)
	//ctx.Abort()
	ResJSON(ctx, em.Status, res)
	return res
}

// ResJSON 响应JSON数据
// 禁止service层调用
func ResJSON(ctx *gin.Context, status int, v interface{}) {
	if ctx == nil {
		return
	}
	buf, err := JSONMarshal(v)
	if err != nil {
		panic(err)
	}
	if status == 0 {
		status = http.StatusOK
	}
	ctx.Data(status, ResponseTypeJSON, buf)
	ctx.Abort()

	// ctx.JSON(status, v)
	// ctx.PureJSON(status, v)
}

// FixResponseError 上级应用已经处理了返回值
func FixResponseError(c *gin.Context, err error) bool {
	switch err.(type) {
	case *Success, *ErrorInfo:
		ResJSON(c, http.StatusOK, err)
		return true
	case *ErrorRedirect:
		status := err.(*ErrorRedirect).Status
		if status <= 0 {
			status = http.StatusSeeOther
		}
		c.Redirect(status, err.(*ErrorRedirect).Location)
		return true
	case *ErrorNone:
		// do nothing
		return true
	case *ErrorModel:
		em := err.(*ErrorModel)
		ResJSON(c, em.Status, NewWrapError(c, em))
		return true
	default:
		// e := err.Error()
		return false
	}
}

// FixResponse400Error 修复返回的异常
func FixResponse400Error(c *gin.Context, err error, errfunc func()) {
	if FixResponseError(c, err) {
		return
	}
	if errfunc != nil {
		errfunc()
	}
	ResError(c, Err400BadRequest)
}

// FixResponse401Error 修复返回的异常, 注意, 401异常会导致系统重定向到登陆页面
func FixResponse401Error(c *gin.Context, err error, errfunc func()) {
	if FixResponseError(c, err) {
		return
	}
	if errfunc != nil {
		errfunc()
	}
	ResError(c, Err401Unauthorized)
}

// FixResponse403Error 修复返回的异常
func FixResponse403Error(c *gin.Context, err error, errfunc func()) {
	if FixResponseError(c, err) {
		return
	}
	if errfunc != nil {
		errfunc()
	}
	ResError(c, Err403Forbidden)
}

// FixResponse406Error 修复返回的异常
func FixResponse406Error(c *gin.Context, err error, errfunc func()) {
	if FixResponseError(c, err) {
		return
	}
	if errfunc != nil {
		errfunc()
	}
	ResError(c, Err406NotAcceptable)
}

// FixResponse500Error 修复返回的异常
func FixResponse500Error(c *gin.Context, err error, errfunc func()) {
	if FixResponseError(c, err) {
		return
	}
	if errfunc != nil {
		errfunc()
	}
	ResError(c, Err500InternalServer)
}

// FixResponse500Error2Logger 修复返回的异常
func FixResponse500Error2Logger(c *gin.Context, err error) {
	FixResponse500Error(c, err, func() { logger.Errorf(c, logger.ErrorWW(err)) })
}

// ResErrorResBody 包装响应错误
// 禁止service层调用
func ResErrorResBody(ctx *gin.Context, em *ErrorModel) error {
	res := NewWrapError(ctx, em)
	ResJSONResBody(ctx, em.Status, res)
	return res
}

// ResJSONResBody 响应JSON数据
// 禁止service层调用
func ResJSONResBody(ctx *gin.Context, status int, v interface{}) {
	if ctx == nil {
		return
	}
	buf, err := JSONMarshal(v)
	if err != nil {
		panic(err)
	}
	ctx.Set(ResBodyKey, buf)
	if status == 0 {
		status = http.StatusOK
	}
	ctx.Data(status, ResponseTypeJSON, buf)
	ctx.Abort()

	// ctx.JSON(status, v)
	// ctx.PureJSON(status, v)
}
