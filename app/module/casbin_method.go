package module

import (
	"time"

	"github.com/suisrc/zgo/modules/helper"
)

// CustomMatchFunc domain
func CustomMatchFunc(args ...interface{}) (interface{}, error) {
	if len(args) != 3 {
		return false, nil
	}
	c8n, b := args[0].(string)
	if !b || c8n == "" {
		return false, nil
	}
	conditions := make(map[string]interface{})
	if err := helper.JSONUnmarshal([]byte(c8n), &conditions); err != nil {
		// panic(err)
		// log.Println(err)
		return false, nil
		// return false, err
		// return false, &helper.ErrorModel{
		// 	Status:   403,
		// 	ShowType: helper.ShowWarn,
		// 	ErrorMessage: &i18n.Message{
		// 		ID:    "ERR-CASBIN-CONDITION",
		// 		Other: "验证器条件错误，拒绝访问",
		// 	},
		// }
	}
	// sub, b := args[1].(CasbinSubject)
	// if !b {
	// 	return false, nil
	// }
	// obj, b := args[2].(CasbinObject)
	// if !b {
	// 	return false, nil
	// }

	result := false
	if access, b := conditions["access_time"]; b {
		if r, e := customAccessTimes(access); e != nil || !r {
			return false, e // 条件失败
		}
		result = true
	}
	return result, nil
}

// 验证授权时间
func customAccessTimes(access interface{}) (bool, error) {
	if access2, b := access.(map[string]interface{}); b {
		if times, b := access2["times"]; b {
			if times2, b := times.([]interface{}); b && len(times2) == 2 {
				now := time.Now()
				if times2[0] != "" {
					if t, e := time.ParseInLocation("2006-01-02 15:04:05", times2[0].(string), time.Local); e != nil {
						return false, nil
					} else if t.After(now) {
						return false, nil
					}
				}
				if times2[1] != "" {
					if t, e := time.ParseInLocation("2006-01-02 15:04:05", times2[1].(string), time.Local); e != nil {
						return false, nil
					} else if t.Before(now) {
						return false, nil
					}
				}
				// log.Println(access)
				return true, nil
			}
		}
	}
	return false, nil
}
