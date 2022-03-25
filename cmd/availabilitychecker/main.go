package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/urfave/cli/v2"

	"github.com/ivan1993spb/availabilitychecker/internal/app"
	"github.com/ivan1993spb/availabilitychecker/internal/core"
	"github.com/ivan1993spb/availabilitychecker/internal/log"
	"github.com/ivan1993spb/availabilitychecker/internal/mysql"
)

var dsn = "checks:checks@tcp(127.0.0.1:3306)/checks?charset=utf8mb4&parseTime=True&loc=Local"

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// TODO: read flags, envs, etc.

	logger := log.NewLogger()
	ctx = log.NewContext(ctx, logger)

	// TODO: configure application.

	app := &cli.App{
		Usage: "availability checker application",
		Commands: []*cli.Command{
			{
				Name:  "mock",
				Usage: "run mocks",
				Action: func(c *cli.Context) error {
					app.NewApp(app.Config{
						Fetcher:  core.RandomFetcher{},
						Executer: core.RandExecutor{},
						Store:    core.Ignore{},
					}).Run(ctx)

					return nil
				},
			},
			{
				Name:  "serve",
				Usage: "run availability service",
				Action: func(c *cli.Context) error {
					db, err := mysql.Connect(dsn)
					if err != nil {
						return err
					}

					err = db.AutoMigrate(&mysql.Check{})
					if err != nil {
						return err
					}

					app.NewApp(app.Config{
						Fetcher: &mysql.Source{
							DB: db,
						},
						Executer: core.CheckExecutor{},
						Store: &mysql.Store{
							DB: db,
						},
					}).Run(ctx)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		logger.Fatal(err)
	}
}
