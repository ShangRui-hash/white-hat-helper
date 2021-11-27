package main

import (
	"os"

	"white-hat-helper/settings"

	"white-hat-helper/controllers"

	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func main() {
	author := cli.Author{
		Name:  "无在无不在",
		Email: "2227627947@qq.com",
	}
	app := &cli.App{
		Name:      "white-hat-helper",
		Usage:     "white-hat-helper",
		UsageText: "white-hat-helper",
		Version:   "v0.1",
		Authors:   []*cli.Author{&author},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "redis",
				Aliases:     []string{"r"},
				Usage:       "redis config file",
				Required:    true,
				Destination: &settings.CurrentConfig.RedisConfigFile,
			},
			&cli.StringFlag{
				Name:        "domains",
				Aliases:     []string{"d"},
				Usage:       "domains such as lenovo.com,lenovo.com.cn,lenovomm.com,lenovo.cn",
				Destination: &settings.CurrentConfig.Domains,
			},
			&cli.StringFlag{
				Name:        "domainFile",
				Aliases:     []string{"f"},
				Usage:       "domain file",
				Destination: &settings.CurrentConfig.DomainFile,
			},
			&cli.BoolFlag{
				Name:        "debug",
				Usage:       "debug",
				Destination: &settings.CurrentConfig.Debug,
			},
			&cli.IntFlag{
				Name:        "companyID",
				Usage:       "companyID",
				Aliases:     []string{"cid"},
				Required:    true,
				Destination: &settings.CurrentConfig.CompanyID,
			},
		},
		Action: controllers.Run,
	}
	err := app.Run(os.Args)
	if err != nil {
		logrus.Error("app.Run faild,err:", err)
	}
}
