package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/request"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/response"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/lender"
)

type LenderHandler struct {
	lenderService *lender.LenderService
	validate      *validator.Validate
}

func NewLenderHandler(lenderService *lender.LenderService, validate *validator.Validate) *LenderHandler {
	return &LenderHandler{
		lenderService: lenderService,
		validate:      validate,
	}
}

// CreateLender godoc
// @Summary Create a new lender
// @Description Register a new lender in the system
// @Tags lenders
// @Accept json
// @Produce json
// @Param lender body request.CreateLenderRequest true "Lender information"
// @Success 201 {object} response.APIResponse{data=lender.Lender} "Lender created successfully"
// @Failure 400 {object} response.APIResponse "Invalid request or validation error"
// @Router /lenders [post]
func (h *LenderHandler) CreateLender(c echo.Context) error {
	var req request.CreateLenderRequest
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

	lender, err := h.lenderService.CreateLender(c.Request().Context(), req.FullName, req.Email, req.PhoneNumber, req.IDNumber)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(lender, "Lender created successfully"))
}

// ListLenders godoc
// @Summary List all lenders
// @Description Get a list of all lenders with optional filtering
// @Tags lenders
// @Accept json
// @Produce json
// @Param full_name query string false "Filter by full name"
// @Param email query string false "Filter by email"
// @Param phone_number query string false "Filter by phone number"
// @Param id_number query string false "Filter by ID number"
// @Success 200 {object} domain.PaginatedResponse{data=[]lender.Lender} "List of lenders"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /lenders [get]
func (h *LenderHandler) ListLenders(c echo.Context) error {
	// Extract query parameters for filtering
	fullName := c.QueryParam("full_name")
	email := c.QueryParam("email")
	phoneNumber := c.QueryParam("phone_number")
	idNumber := c.QueryParam("id_number")

	filter := lender.LenderFilter{
		FullName:    &fullName,
		Email:       &email,
		PhoneNumber: &phoneNumber,
		IDNumber:    &idNumber,
	}

	lenders, err := h.lenderService.ListLenders(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, lenders)
}
