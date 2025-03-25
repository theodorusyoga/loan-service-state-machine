package main

import (
	"github.com/theodorusyoga/loan-service-state-machine/config"
	fxpkg "github.com/theodorusyoga/loan-service-state-machine/pkg/fx"
	"go.uber.org/fx"

	_ "github.com/theodorusyoga/loan-service-state-machine/docs"
)

// @title Loan Service API
// @version 1.0
// @description API for managing loans with a state machine workflow
// @host localhost:5002
// @BasePath /api/v1
func main() {
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
