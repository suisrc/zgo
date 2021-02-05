package module

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/modules/auth"
	"github.com/suisrc/zgo/modules/helper"

	"github.com/gin-gonic/gin"
)

// QueryServiceCode 查询服务
// "zgo:svc-cox:" + host + ":" + resource
func (a *CasbinAuther) QueryServiceCode(ctx *gin.Context, user auth.UserInfo, host, path, org string) (string, int64, error) {
	resource := ""
	if strings.HasPrefix(path, "/api/") {
		// 后端API服务使用3级模糊匹配
		resource = "/" + helper.SplitStrCR(path[1:], '/', 3)
	}
	if host == "" && resource == "" {
		return "", 0, errors.New("no service")
	}
	// audience := helper.ReverseStr(host) // host倒序
	audience := host
	key := "zgo:svc-cox:" + audience + ":" + resource

	if svc, b, err := a.Storer.Get(ctx, key); err != nil {
		return "", 0, err // 查询缓存出现异常
	} else if b {
		if strings.HasPrefix(svc, "err:") {
			return "", 0, errors.New(svc[4:]) // 上一次查询，拒绝请求
		}
		offset := strings.IndexRune(svc, '/')
		if offset <= 0 {
			a.Storer.Delete(ctx, key)
			return "", 0, errors.New("系统缓存异常:[" + key + "]" + svc)
		}
		sid, _ := strconv.Atoi(svc[offset+1:])
		return svc[:offset], int64(sid), nil
	}

	// 由于查询是居于全局的， 所以1分钟的缓存是一个合理的范围
	sa := schema.CasbinGpaSvcAud{}
	if err := sa.QueryByAudAndResAndOrg(a.Sqlx, audience, resource, ""); err != nil && !sqlxc.IsNotFound(err) {
		// 系统没有配置或者系统为指定有效服务名称
		a.Storer.Set(ctx, key, "err:"+err.Error(), CasbinServiceCodeExpireAt) // 1分钟延迟刷新， 拒绝请求也需要缓存
		return "", 0, err
	} else if !sa.SvcCode.Valid {
		a.Storer.Set(ctx, key, "err:no service", CasbinServiceCodeExpireAt)
		return "", 0, errors.New("no service")
	}
	a.Storer.Set(ctx, key, sa.SvcCode.String+"/"+strconv.Itoa(int(sa.SvcID.Int64)), CasbinServiceCodeExpireAt) // 查询结果缓存1分钟
	return sa.SvcCode.String, sa.SvcID.Int64, nil
}

// CheckTenantService 验证租户是否有访问该服务的权限服务
// "zgo:svc-orx:" + svc_cod + ":" + org_cod -> CasbinGpaSvcOrg
func (a *CasbinAuther) CheckTenantService(ctx *gin.Context, user auth.UserInfo, org, svc string, sid int64) (bool, error) {
	if org == "" || org == schema.PlatformCode {
		return true, nil // 平台用户， 没有服务权限问题
	}

	key := "zgo:svc-orx:" + svc + ":" + org
	if res, b, err := a.Storer.Get(ctx, key); err != nil {
		return false, err
	} else if b {
		if res == "1" {
			return true, nil
		}
		offset := strings.IndexRune(res, '/')
		if offset <= 0 {
			a.Storer.Delete(ctx, key)
			return false, errors.New("系统缓存异常:[" + key + "]" + res)
		}
		return false, helper.New0Error(ctx, helper.ShowWarn, &i18n.Message{ID: res[:offset], Other: res[offset+1:]})
	}

	var emsg *i18n.Message
	so := schema.CasbinGpaSvcOrg{}
	// 1:启用 0:禁用 2:未激活 3: 删除 4: 欠费 5: 到期
	if err := so.QueryByOrgAndSvc(a.Sqlx, org, sid); err != nil {
		if !sqlxc.IsNotFound(err) {
			return false, err // 系统内部的位置异常
		}
		emsg = &i18n.Message{ID: "WARN-SERVICE-NOFOUND", Other: "访问的服务不存在"}
	} else if so.Expired.Valid && time.Now().After(so.Expired.Time) {
		// 前置授权异常
		emsg = &i18n.Message{ID: "WARN-SERVICE-EXPIRED", Other: "授权已经过期"}
	} else if so.Status == schema.StatusEnable {
		// 正常结果返回
		expiration := CasbinServiceTenantExpireAt // 延迟刷新
		if so.Expired.Valid && so.Expired.Time.Sub(time.Now()) < expiration {
			expiration = so.Expired.Time.Sub(time.Now()) // 修正过期时间
		}
		a.Storer.Set(ctx, key, "1", expiration)
		return true, nil
	} else if so.Status == schema.StatusDisable {
		emsg = &i18n.Message{ID: "WARN-SERVICE-DISABLE", Other: "服务已经被禁用"}
	} else if so.Status == schema.StatusDeleted {
		emsg = &i18n.Message{ID: "WARN-SERVICE-DELETE", Other: "服务已经被删除"}
	} else if so.Status == schema.StatusNoActivate {
		emsg = &i18n.Message{ID: "WARN-SERVICE-NOACTIVATE", Other: "服务未激活"}
	} else if so.Status == schema.StatusExpired {
		emsg = &i18n.Message{ID: "WARN-SERVICE-EXPIRED", Other: "授权已经过期"}
	} else {
		emsg = &i18n.Message{ID: "WARN-SERVICE-OTHER", Other: "授权状态异常"}
	}
	a.Storer.Set(ctx, key, emsg.ID+"/"+emsg.Other, CasbinServiceTenantExpireAt/4) // 拒绝请求也需要缓存, 时间缩短1/4
	return false, helper.New0Error(ctx, helper.ShowWarn, emsg)
}
