// Code generated by Wire. DO NOT EDIT.

//go:generate wire
//+build !wireinject

package injector

import (
	"github.com/suisrc/zgo/app/api"
	"github.com/suisrc/zgo/app/api/manager"
	"github.com/suisrc/zgo/app/model/entc"
	"github.com/suisrc/zgo/app/model/gpa"
	"github.com/suisrc/zgo/app/model/sqlxc"
	"github.com/suisrc/zgo/app/module"
	"github.com/suisrc/zgo/app/oauth2"
	"github.com/suisrc/zgo/app/service"
	"github.com/suisrc/zgo/middlewire"
	"github.com/suisrc/zgo/modules/passwd"
)

// Injectors from wire.go:

func BuildInjector() (*Injector, func(), error) {
	bundle := NewBundle()
	useEngine := api.NewUseEngine(bundle)
	engine := middlewire.InitGinEngine(useEngine)
	router := middlewire.NewRouter(engine)
	client, cleanup, err := entc.NewClient()
	if err != nil {
		return nil, nil, err
	}
	db, cleanup2, err := sqlxc.NewClient()
	if err != nil {
		cleanup()
		return nil, nil, err
	}
	gpaGPA := gpa.GPA{
		Entc: client,
		Sqlx: db,
	}
	storer, cleanup3, err := module.NewStorer()
	if err != nil {
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	authOpts := &module.AuthOpts{
		GPA:    gpaGPA,
		Storer: storer,
	}
	auther := module.NewAuther(authOpts)
	casbinAuther := &module.CasbinAuther{
		GPA:    gpaGPA,
		Storer: storer,
		Auther: auther,
	}
	auth := &api.Auth{
		CasbinAuther: casbinAuther,
	}
	validator := &passwd.Validator{}
	mobileSender := service.MobileSender{
		GPA: gpaGPA,
	}
	emailSender := service.EmailSender{
		GPA: gpaGPA,
	}
	threeSender := service.ThreeSender{
		GPA: gpaGPA,
	}
	selector, err := oauth2.NewSelector(gpaGPA, storer)
	if err != nil {
		cleanup3()
		cleanup2()
		cleanup()
		return nil, nil, err
	}
	signin := service.Signin{
		GPA:            gpaGPA,
		Passwd:         validator,
		Store:          storer,
		MSender:        mobileSender,
		ESender:        emailSender,
		TSender:        threeSender,
		OAuth2Selector: selector,
	}
	apiSignin := &api.Signin{
		GPA:           gpaGPA,
		Auther:        auther,
		SigninService: signin,
		CasbinAuther:  casbinAuther,
	}
	connect := &api.Connect{
		GPA:          gpaGPA,
		Signin:       apiSignin,
		CasbinAuther: casbinAuther,
	}
	user := service.User{
		GPA:            gpaGPA,
		Store:          storer,
		OAuth2Selector: selector,
	}
	apiUser := &api.User{
		GPA:          gpaGPA,
		UserService:  user,
		CasbinAuther: casbinAuther,
	}
	system := &api.System{
		GPA: gpaGPA,
	}
	weixin := &api.Weixin{
		GPA: gpaGPA,
	}
	managerUser := &manager.User{
		GPA: gpaGPA,
	}
	account := &manager.Account{
		GPA: gpaGPA,
	}
	role := &manager.Role{
		GPA: gpaGPA,
	}
	menu := &manager.Menu{
		GPA: gpaGPA,
	}
	gateway := &manager.Gateway{
		GPA: gpaGPA,
	}
	wire := &manager.Wire{
		User:    managerUser,
		Account: account,
		Role:    role,
		Menu:    menu,
		Gateway: gateway,
	}
	options := &api.Options{
		Engine:       engine,
		Router:       router,
		Auth:         auth,
		Connect:      connect,
		Signin:       apiSignin,
		User:         apiUser,
		System:       system,
		Weixin:       weixin,
		CasbinAuther: casbinAuther,
		ManagerWire:  wire,
	}
	endpoints := api.InitEndpoints(options)
	i18n := &service.I18n{
		GPA:    gpaGPA,
		Bundle: bundle,
	}
	i18nLoader := service.InitI18nLoader(i18n)
	swagger := middlewire.NewSwagger(engine)
	healthz := middlewire.NewHealthz(engine)
	injector := &Injector{
		Engine:     engine,
		Endpoints:  endpoints,
		Bundle:     bundle,
		I18nLoader: i18nLoader,
		Swagger:    swagger,
		Healthz:    healthz,
	}
	return injector, func() {
		cleanup3()
		cleanup2()
		cleanup()
	}, nil
}
