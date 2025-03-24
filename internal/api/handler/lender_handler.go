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
