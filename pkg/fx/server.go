package fx

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/theodorusyoga/loan-service-state-machine/config"
	"go.uber.org/fx"
)

const ServiceName = "loan-service"

type ServerParams struct {
	fx.In

	Lifecycle fx.Lifecycle
	Config    *config.Config
}

func NewServer(p ServerParams) *echo.Echo {
	e := echo.New()

	// TODO: Middleware and routing setup can be done here

	p.Lifecycle.Append(fx.Hook{
		OnStart: func(context.Context) error {
			serverAddr := fmt.Sprintf(":%s", p.Config.Server.Port)
			go func() {
				e.Logger.Infof("Starting server on %s", serverAddr)
				if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
					e.Logger.Fatal("Shutting down the %s service: ", ServiceName, err)
				}
			}()
			return nil
		},
		OnStop: func(context.Context) error {
			return e.Shutdown(context.Background())
		},
	})

	return e
}
