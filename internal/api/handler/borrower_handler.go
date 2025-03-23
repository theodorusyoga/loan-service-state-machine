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
