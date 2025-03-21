package fx

import (
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
	"github.com/theodorusyoga/loan-service-state-machine/internal/repository"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewServer,
	),
	DomainModule,
	InfrastructureModule,
)

var DomainModule = fx.Module("domain", fx.Provide(
	loan.NewDefaultStatusValidator,
	loan.NewLoanService,
))

var InfrastructureModule = fx.Module("infrastructure",
	fx.Provide(
		// Database
		repository.NewDatabase,

		// Repositories
		fx.Annotate(
			repository.NewLoanRepository,
			fx.As(new(loan.Repository)),
		),
	),
)
