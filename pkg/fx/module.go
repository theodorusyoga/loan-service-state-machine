package fx

import (
	"github.com/theodorusyoga/loan-service-state-machine/config"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewServer,
		config.Load,
	),
)
