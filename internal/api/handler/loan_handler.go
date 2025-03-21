package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/request"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/response"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

type LoanHandler struct {
	loanService *loan.LoanService
	validate    *validator.Validate
}

func NewLoanHandler(loanService *loan.LoanService, validate *validator.Validate) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
		validate:    validate,
	}
}

func (h *LoanHandler) CreateLoan(c echo.Context) error {
	var req request.CreateLoanRequest
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

	loan, err := h.loanService.CreateLoan(c.Request().Context(), req.BorrowerID, req.Amount, req.Rate, req.ROI)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(loan, "Loan created successfully"))
}

func (h *LoanHandler) ApproveLoan(c echo.Context) error {
	id := c.Param("id")

	var req request.ApprovalRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request"))
	}

	loan := &loan.Loan{
		ID:     id,
		Status: loan.StatusProposed,
	}

	err := h.loanService.ApproveLoan(loan, req.ApprovedBy)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("Loan approved successfully"))
}

func formatValidationErrors(errors validator.ValidationErrors) string {
	var errorMsg string
	for _, err := range errors {
		switch err.Tag() {
		case "required":
			errorMsg += err.Field() + " is required. "
		case "gt":
			errorMsg += err.Field() + " must be greater than " + err.Param() + ". "
		case "lt":
			errorMsg += err.Field() + " must be less than " + err.Param() + ". "
		case "roiLessThanRate":
			errorMsg += "ROI must be less than Rate. "
		default:
			errorMsg += err.Field() + " is invalid. "
		}
	}
	return errorMsg
}
