package main

import (
	"github.com/theodorusyoga/loan-service-state-machine/config"
	fxpkg "github.com/theodorusyoga/loan-service-state-machine/pkg/fx"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		fx.Provide(
			func() string {
				return "config/development.yaml"
			},
			config.Load,
		),
		fxpkg.Module,
		// TODO: Put other modules
	)

	app.Run()
}
