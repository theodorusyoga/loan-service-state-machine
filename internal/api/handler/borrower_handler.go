package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/request"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/response"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/borrower"
)

type BorrowerHandler struct {
	borrowerService *borrower.BorrowerService
	validate        *validator.Validate
}

func NewBorrowerHandler(borrowerService *borrower.BorrowerService, validate *validator.Validate) *BorrowerHandler {
	return &BorrowerHandler{
		borrowerService: borrowerService,
		validate:        validate,
	}
}

// CreateBorrower godoc
// @Summary Create a new borrower
// @Description Register a new borrower in the system
// @Tags borrowers
// @Accept json
// @Produce json
// @Param borrower body request.CreateBorrowerRequest true "Borrower information"
// @Success 201 {object} response.APIResponse{data=borrower.Borrower} "Borrower created successfully"
// @Failure 400 {object} response.APIResponse "Invalid request or validation error"
// @Router /borrowers [post]
func (h *BorrowerHandler) CreateBorrower(c echo.Context) error {
	var req request.CreateBorrowerRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request"))
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		// Format validation errors nicely
		validationErrors := err.(validator.ValidationErrors)
		errorsMsg := formatValidationErrors(validationErrors)
		return c.JSON(http.StatusBadRequest, response.Error(errorsMsg))
	}

	borrower, err := h.borrowerService.CreateBorrower(c.Request().Context(), req.FullName, req.Email, req.PhoneNumber, req.IDNumber)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(borrower, "Borrower created successfully"))
}

// ListBorrowers godoc
// @Summary List all borrowers
// @Description Get a list of all borrowers with optional filtering
// @Tags borrowers
// @Accept json
// @Produce json
// @Param full_name query string false "Filter by full name"
// @Param email query string false "Filter by email"
// @Param phone_number query string false "Filter by phone number"
// @Param id_number query string false "Filter by ID number"
// @Success 200 {object} domain.PaginatedResponse{data=[]borrower.Borrower} "List of borrowers"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /borrowers [get]
func (h *BorrowerHandler) ListBorrowers(c echo.Context) error {
	// Extract query parameters for filtering
	fullName := c.QueryParam("full_name")
	email := c.QueryParam("email")
	phoneNumber := c.QueryParam("phone_number")
	idNumber := c.QueryParam("id_number")

	filter := borrower.BorrowerFilter{
		FullName:    &fullName,
		Email:       &email,
		PhoneNumber: &phoneNumber,
		IDNumber:    &idNumber,
	}

	borrowers, err := h.borrowerService.ListBorrowers(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, borrowers)
}
