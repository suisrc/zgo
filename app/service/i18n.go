package service

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/schema"
	"github.com/suisrc/zgo/modules/config"
	"github.com/suisrc/zgo/modules/logger"
)

// InitI18nLoader handler
func InitI18nLoader(loader *I18n) I18nLoader {
	if config.C.I18N.DBEnable {
		loader.LoadI18nMessage() // 加载数据库的配置
		return loader
	}
	return nil
}

// I18nLoader loader
type I18nLoader *I18n

// I18n 账户管理
type I18n struct {
	gpa.GPA              // 数据库
	Bundle  *i18n.Bundle // 控制器
}

// LoadI18nMessage load
func (a *I18n) LoadI18nMessage() error {
	count := 0
	if i18ns, err := new(schema.I18nGpaMessage).QueryAll(a.Sqlx); err != nil {
		logger.Errorf(nil, "i18n add message error: %s", logger.ErrorWW(err))
		return err
	} else if i18ns != nil {
		for _, m := range *i18ns {
			message, language := m.I18nMessage()
			if err := a.Bundle.AddMessages(language, message); err != nil {
				logger.Errorf(nil, "i18n add message error: %s", logger.ErrorWW(err))
			} else {
				count++
			}
		}
	}
	logger.Infof(nil, "i18n load message from database, count: %d", count)
	return nil
}
