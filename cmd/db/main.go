package main

import (
	"context"
	"log"
	"os"

	"github.com/suisrc/zgo/cmd/db/mysql"

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
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Println(err.Error())
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
			mf := &mysql.ModelFile{
				Model:  c.String("model"),
				Output: c.String("output"),
			}
			return mf.RunBuild()
		},
	}
}
