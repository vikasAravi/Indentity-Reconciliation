package main

import (
	"bitespeed/identity-reconciliation/appcontext"
	"bitespeed/identity-reconciliation/config"
	"bitespeed/identity-reconciliation/services"
	"bitespeed/identity-reconciliation/utils"
	"context"
	"fmt"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	commandArgs := os.Args
	config.SetConfigFileFromArgs(commandArgs)
	config.Load()

	util.SetupLogger()
	ctx := context.Background()
	err := appcontext.Initiate(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	app := &cli.App{
		Name:  "Reconciliation Service",
		Usage: "Reconciliation Service",
		Commands: []*cli.Command{
			{
				Name:   "api",
				Usage:  "run api",
				Action: services.StartAPI,
			},
		},
	}

	if e := app.Run(os.Args); e != nil {
		fmt.Println(e)
	}
}

// 1. create model schemas
// 2. migration files
// 3. API - schemas
// 4. Query and Persistence finalization's
// 5. Dockerfile
