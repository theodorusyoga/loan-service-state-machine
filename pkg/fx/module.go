package fx

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/theodorusyoga/loan-service-state-machine/config"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/handler"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/borrower"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
	"github.com/theodorusyoga/loan-service-state-machine/internal/repository"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

func ProvideValidator() *validator.Validate {
	v := validator.New()

	// Register a validator for ROI < Rate
	v.RegisterValidation("roiLessThanRate", func(fl validator.FieldLevel) bool {
		// Get the struct
		parent := fl.Parent()

		rateField := parent.FieldByName("Rate")
		if !rateField.IsValid() {
			return false
		}
		roi := fl.Field().Float()
		rate := rateField.Float()

		// ROI must be less than Rate
		return roi < rate
	})

	return v
}

var Module = fx.Options(
	InfrastructureModule,
	DomainModule,
	APIModule,
)

var DomainModule = fx.Module("domain", fx.Provide(
	loan.NewDefaultStatusValidator,
	loan.NewLoanService,
	borrower.NewBorrowerService,
))

var InfrastructureModule = fx.Module("infrastructure",
	fx.Provide(
		// Database
		repository.NewDatabase,

		func(db *repository.Database) *gorm.DB {
			return db.DB
		},

		// Repositories
		fx.Annotate(
			repository.NewLoanRepository,
			fx.As(new(loan.Repository)),
		),
		fx.Annotate(
			repository.NewBorrowerRepository,
			fx.As(new(borrower.Repository)),
		),
	),
)

var APIModule = fx.Module("api", fx.Provide(
	ProvideValidator,
	handler.NewLoanHandler,
	handler.NewBorrowerHandler,
	NewServer,
),
	fx.Invoke(registerRoutes))

func registerRoutes(lc fx.Lifecycle, e *echo.Echo, cfg *config.Config, loanHandler *handler.LoanHandler, borrowerHandler *handler.BorrowerHandler) {
	api := e.Group("/api/v1")

	loans := api.Group("/loans")
	loans.POST("", loanHandler.CreateLoan)
	loans.POST("/:id/approve", loanHandler.ApproveLoan)

	borrowers := api.Group("/borrowers")
	borrowers.POST("", borrowerHandler.CreateBorrower)

	// Start server in a goroutine
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				// Log server startup
				e.Logger.Info("Starting HTTP server on port " + cfg.Server.Port)

				// Log registered routes
				for _, route := range e.Routes() {
					e.Logger.Info("Route registered: " + route.Method + " " + route.Path)
				}

				if err := e.Start(":" + cfg.Server.Port); err != nil {
					if err != http.ErrServerClosed {
						e.Logger.Fatal("Server error: " + err.Error())
					} else {
						e.Logger.Info("Server stopped")
					}
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return e.Shutdown(ctx)
		},
	})
}

func NewServer(cfg *config.Config) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	return e
}
