package service

import (
	"github.com/nicksnyder/go-i18n/v2/i18n"
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
	GPA                 // 数据库
	Bundle *i18n.Bundle // 控制器
}

// LoadI18nMessage load
func (a *I18n) LoadI18nMessage() error {
	i18n0 := schema.I18nGpaMessage{}
	i18ns := []schema.I18nGpaMessage{}
	err := a.Sqlx.Select(&i18ns, i18n0.SQLByALL())
	if err != nil {
		logger.Errorf(nil, "i18n add message error: %s", err.Error())
		return err
	}
	count := 0
	for _, m := range i18ns {
		message, language := m.I18nMessage()
		err2 := a.Bundle.AddMessages(language, message)
		if err2 != nil {
			logger.Errorf(nil, "i18n add message error: %s", err2.Error())
		} else {
			count++
		}
	}
	logger.Infof(nil, "i18n load message from database, count: %d", count)
	return nil
}
