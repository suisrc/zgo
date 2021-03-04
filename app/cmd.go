package app

import (
	"context"

	"github.com/urfave/cli/v2"
)

// SetConfigFile 设定配置文件
func SetConfigFile(s string) Option {
	return func(o *Options) {
		o.ConfigFile = s
	}
}

// SetVersion 设定版本号b
func SetVersion(s string) Option {
	return func(o *Options) {
		o.Version = s
	}
}

// SetBuildInjector 设定注入助手
func SetBuildInjector(b BuildInjector) Option {
	return func(o *Options) {
		o.BuildInjector = b
	}
}

// RunWebCmd ...
func RunWebCmd(ctx context.Context, action func(c *cli.Context) error) *cli.Command {
	return &cli.Command{
		Name:  "web",
		Usage: "运行web服务",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "conf",
				Aliases:     []string{"c"},
				Usage:       "配置文件(.json,.yaml,.toml)",
				DefaultText: "config.toml",
				//Required:   true,
			},
		},
		Action: action,
	}
}
