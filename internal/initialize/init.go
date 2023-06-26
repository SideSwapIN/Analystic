package initialize

import (
	"os"

	"github.com/SideSwapIN/Analystic/internal/config"
	"github.com/SideSwapIN/Analystic/internal/db"
	"github.com/urfave/cli/v2"
)

func Init() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"conf"},
				Value:   "config.toml",
				Usage:   "configure toml file path",
			},
		},
		Action: func(ctx *cli.Context) error {
			err := config.Init(ctx.String("config"))
			if err != nil {
				return err
			}

			if err := db.InitMysqlDB(); err != nil {
				return err
			}

			if err := db.InitRedisDB(); err != nil {
				return err
			}

			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
