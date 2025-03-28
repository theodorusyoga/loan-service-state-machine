package fx

import (
	"context"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/theodorusyoga/loan-service-state-machine/config"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/handler"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/borrower"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/document"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/employee"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/lender"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan/callbacks"
	loanlender "github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan_lender"
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
	employee.NewEmployeeService,
	document.NewDocumentService,
	lender.NewLenderService,
	loanlender.NewLoanLenderService,

	// Callback registrar for FSM
	fx.Annotate(
		callbacks.New,
		fx.As(new(loan.CallbackRegistrar)),
	),
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
		fx.Annotate(
			repository.NewEmployeeRepository,
			fx.As(new(employee.Repository)),
		),
		fx.Annotate(
			repository.NewDocumentRepository,
			fx.As(new(document.Repository)),
		),
		fx.Annotate(
			repository.NewLenderRepository,
			fx.As(new(lender.Repository)),
		),
		fx.Annotate(
			repository.NewLoanLenderRepository,
			fx.As(new(loanlender.Repository)),
		),
	),
)

var APIModule = fx.Module("api", fx.Provide(
	ProvideValidator,
	handler.NewLoanHandler,
	handler.NewBorrowerHandler,
	handler.NewEmployeeHandler,
	handler.NewLenderHandler,
	NewServer,
),
	fx.Invoke(registerRoutes))

func registerRoutes(lc fx.Lifecycle,
	e *echo.Echo, cfg *config.Config, loanHandler *handler.LoanHandler,
	borrowerHandler *handler.BorrowerHandler, emp *handler.EmployeeHandler,
	lenderHandler *handler.LenderHandler) {
	api := e.Group("/api/v1")

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	loans := api.Group("/loans")
	loans.GET("", loanHandler.ListLoans)
	loans.POST("", loanHandler.CreateLoan)
	loans.PATCH("/:id/:status", loanHandler.UpdateLoanStatus)

	borrowers := api.Group("/borrowers")
	borrowers.GET("", borrowerHandler.ListBorrowers)
	borrowers.POST("", borrowerHandler.CreateBorrower)

	employees := api.Group("/employees")
	employees.GET("", emp.ListEmployees)
	employees.POST("", emp.CreateEmployee)

	lenders := api.Group("/lenders")
	lenders.GET("", lenderHandler.ListLenders)
	lenders.POST("", lenderHandler.CreateLender)

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
