package main

import (
	"context"
	"log"
	"os"

	entcmd "github.com/suisrc/zgo/cmd/db/ent"
	mysqlcmd "github.com/suisrc/zgo/cmd/db/mysql"
	"github.com/suisrc/zgo/modules/logger"

	"github.com/urfave/cli/v2"
)

// VERSION 版本号，可以通过编译的方式指定版本号：go build -ldflags "-X main.VERSION=x.x.x"
var VERSION = "0.0.1"

func main() {
	ctx := context.Background()

	app := cli.NewApp()
	app.Name = "zgo"
	app.Version = VERSION
	app.Usage = "zgo cmd"
	app.Commands = []*cli.Command{
		runMysqlCmd(ctx),
		runEntcCmd(ctx),
		runEntcDel(ctx),
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Println(logger.ErrorWW(err))
	}
}

func runMysqlCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "mysql",
		Usage: "构建mysql文",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "model",
				Aliases:     []string{"m"},
				Usage:       "输入文件(.md)",
				DefaultText: "doc/model.md",
				//Required:   true,
			},
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Usage:       "输出文件(.sql)",
				DefaultText: "doc/model.sql",
				//Required:   true,
			},
		},
		Action: func(c *cli.Context) error {
			mf := &mysqlcmd.ModelFile{
				Model:  c.String("model"),
				Output: c.String("output"),
			}
			return mf.RunBuild()
		},
	}
}

func runEntcCmd(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "entc",
		Usage: "构建ent文",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "model",
				Aliases: []string{"m"},
				Usage:   "输入文件(.md)",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Usage:   "输出文件文件夹",
			},
		},
		Action: func(c *cli.Context) error {
			mf := &entcmd.ModelFile{
				Model:  c.String("model"),
				Output: c.String("output"),
			}
			return mf.RunBuild()
		},
	}
}

func runEntcDel(ctx context.Context) *cli.Command {
	return &cli.Command{
		Name:  "entc-del",
		Usage: "删除ent中的实体",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "entitys",
				Aliases: []string{"es"},
				Usage:   "删除实体列表",
			},
			&cli.StringFlag{
				Name:    "input",
				Aliases: []string{"o"},
				Usage:   "输入文件夹",
			},
		},
		Action: func(c *cli.Context) error {
			return entcmd.RunDel(c.String("input"), c.String("entitys"))
		},
	}
}
