package module

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/modules/helper"
)

// CasbinPolicy Casbin策略
type CasbinPolicy struct {
	Mid       int64
	Ver       string
	New       bool       // 重新构建
	ModelText string     // 模型声明
	Grouping  [][]string // 角色声明
	Policies  [][]string // 策略声明
	Version   string     // 策略版本
}

// QueryCasbinPolicies 获取Casbin策略
func (a *CasbinAuther) QueryCasbinPolicies(org, ver string) (*CasbinPolicy, error) {
	c := CasbinPolicy{
		Grouping: [][]string{},
		Policies: [][]string{},
	}
	// 获取策略模型
	cgm := schema.CasbinGpaModel{}
	if err := cgm.QueryByOrg(a.Sqlx, org); err != nil && !sqlxc.IsNotFound(err) {
		// 数据库异常
		return nil, err
	}
	if cgm.ID == 0 {
		// 新建访问策略
		cgm.Name = sql.NullString{Valid: true, String: "Default"}
		cgm.Ver = sql.NullString{Valid: true, String: "1.0.0"}
		cgm.Org = sql.NullString{Valid: true, String: org}
		cgm.Description = sql.NullString{Valid: true, String: "Auto Build"}
		cgm.Status = schema.StatusNoActivate // 未激活状态
		// cgm.Statement = sql.NullString{Valid: true, String: CasbinDefaultMatcher}
		if err := cgm.SaveOrUpdate(a.Sqlx); err != nil {
			return nil, err
		}
	}
	nver := fmt.Sprintf("%s:%s", strconv.Itoa(int(cgm.ID)), cgm.Ver.String)
	if ver != "" && ver == nver {
		return nil, nil
	}
	// 访问策略更新
	c.Mid = cgm.ID
	c.Ver = cgm.Ver.String
	c.Version = fmt.Sprintf("%s:%s", strconv.Itoa(int(cgm.ID)), cgm.Ver.String)
	if cgm.Statement.Valid {
		c.ModelText = CasbinPolicyModel + cgm.Statement.String
	} else {
		c.ModelText = CasbinPolicyModel + CasbinDefaultMatcher
	}
	if cgm.Status == schema.StatusEnable && cgm.HasRules(a.Sqlx2, c.Mid, c.Ver) {
		// 访问策略已经构建完成，不用重新构建
		return &c, nil
	} else if cgm.Status == schema.StatusDisable {
		return nil, &helper.ErrorModel{
			Status:   403,
			ShowType: helper.ShowWarn,
			ErrorMessage: &i18n.Message{
				ID:    "ERR-CASBIN-DISABLE",
				Other: "授权系统已经被禁止使用，请联系平台管理员",
			},
		}
	}

	// 获取基础配置访问策略
	if err := a.CreateCasbinPolicy(org, &c); err != nil {
		return nil, err
	}
	if len(c.Grouping) == 0 && len(c.Policies) == 0 {
		// 无法处理规则, 给出默认无用规则
		// sub,  svc, org, path, meth, eft, c8n
		pp := []string{"none", "", "", "", "", "deny", ""}
		c.Policies = append(c.Policies, pp)
	}
	c.New = true // 模型需要重新构建
	return &c, nil
}

// CreateCasbinPolicy 获取Casbin策略
func (a *CasbinAuther) CreateCasbinPolicy(org string, c *CasbinPolicy) error {
	// log.Println(c.ModelText)
	// 获取角色间的关系
	if rrs, err := new(schema.CasbinGpaRoleRole).QueryByOrg(a.Sqlx, org); err != nil {
		if !sqlxc.IsNotFound(err) {
			return err
		}
		// 没有有效的角色关系
	} else if len(*rrs) > 0 {
		for _, v := range *rrs {
			rr := []string{v.ParentName, v.ChildName}

			// 角色前增加应用标识， 标记应用专有角色
			if v.ParentSvc.Valid {
				rr[0] = v.ParentSvc.String + ":" + rr[0]
			}
			if v.ChildSvc.Valid {
				rr[1] = v.ChildSvc.String + ":" + rr[1]
			}
			// 角色前增加Casbin角色专有前缀
			rr[0] = CasbinRolePrefix + rr[0]
			rr[1] = CasbinRolePrefix + rr[1]
			c.Grouping = append(c.Grouping, rr)
		}
	}
	if rps, err := new(schema.CasbinGpaRolePolicy).QueryByOrg(a.Sqlx, org); err != nil {
		if !sqlxc.IsNotFound(err) {
			return err
		}
		// 没有有效角色策略关系
	} else if len(*rps) > 0 {
		for _, v := range *rps {
			rp := []string{v.Role, v.Policy}

			// 角色前增加应用标识， 标记应用专有角色
			if v.Svc.Valid {
				rp[0] = v.Svc.String + ":" + rp[0]
			}
			// 角色前增加Casbin角色专有前缀
			rp[0] = CasbinRolePrefix + rp[0]
			// 策略前增加Casbin策略专有前缀
			rp[1] = CasbinPolicyPrefix + rp[1]
			c.Grouping = append(c.Grouping, rp)
		}
	}
	if pss, err := new(schema.CasbinGpaPolicyStatement).QueryByOrg(a.Sqlx, org); err != nil {
		if !sqlxc.IsNotFound(err) {
			return err
		}
	} else if len(*pss) > 0 {
		for _, v := range *pss {
			// 策略前增加Casbin策略专有前缀
			sub := CasbinPolicyPrefix + v.Name
			eft := helper.IfString(v.Effect, "allow", "deny")
			c8n := v.Condition.String
			if v.Action.Valid {
				actions := strings.Split(v.Action.String, ";")
				for _, action := range actions {
					sa := strings.SplitN(action, ":", 2)
					if len(sa) != 2 {
						break
					}
					svc := sa[0]
					if pas, err := new(schema.CasbinGpaPolicyServiceAction).QueryActionByNameAndSvc(a.Sqlx, sa[1], sa[0]); err != nil {
						if !sqlxc.IsNotFound(err) {
							return err
						}
					} else if len(*pas) > 0 {
						for _, a := range *pas {
							if a.Resource.Valid {
								paths := strings.Split(a.Resource.String, ";")
								for _, path := range paths {
									meth := "*"
									if offset := strings.IndexRune(path, ' '); offset > 0 {
										meth = path[:offset]
										path = path[offset+1:]
									}
									pp := []string{sub, svc, org, path, meth, eft, c8n}
									c.Policies = append(c.Policies, pp)
								}
							}
						}
					}
				}
			} else if v.Resource.Valid {
				// 配置资源访问权限， 暂时没有进行开发
			}
		}
	}

	return nil
}
