package main

import (
	"github.com/theodorusyoga/loan-service-state-machine/config"
	"github.com/theodorusyoga/loan-service-state-machine/migrations"
	fxpkg "github.com/theodorusyoga/loan-service-state-machine/pkg/fx"
	"go.uber.org/fx"
)

func main() {
	// Running migrations
	migrations.RunMigrations()

	app := fx.New(
		fx.Provide(
			func() string {
				return "config/config.yaml"
			},
			config.Load,
		),
		fxpkg.Module,
		// TODO: Put other modules
	)

	app.Run()
}
