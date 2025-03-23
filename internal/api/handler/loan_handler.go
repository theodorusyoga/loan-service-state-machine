package handler

import (
	"net/http"
	"strconv"

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

func (h *LoanHandler) ListLoans(c echo.Context) error {
	// Extract query parameters for filtering
	maxAmountStr := c.QueryParam("max_amount")
	minAmountStr := c.QueryParam("min_amount")

	var maxAmount, minAmount *float64

	if maxAmountStr != "" {
		val, err := strconv.ParseFloat(maxAmountStr, 64)
		if err == nil {
			maxAmount = &val
		}
	}

	if minAmountStr != "" {
		val, err := strconv.ParseFloat(minAmountStr, 64)
		if err == nil {
			minAmount = &val
		}
	}

	filter := loan.LoanFilter{
		MaxAmount: maxAmount,
		MinAmount: minAmount,
	}

	borrowers, err := h.loanService.ListLoans(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, borrowers)
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

func (h *LoanHandler) UpdateLoanStatus(c echo.Context) error {
	loanID := c.Param("id")
	newStatus := c.Param("status")

	// Validate the requested status
	if !loan.IsValidStatus(newStatus) {
		return c.JSON(http.StatusBadRequest, response.Error("invalid status"))
	}

	var req request.StatusUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("invalid request"))
	}

	// Get loan with ID
	loanEntity, entityErr := h.loanService.GetByID(c.Request().Context(), loanID)
	if entityErr != nil {
		return c.JSON(http.StatusBadRequest, response.Error(entityErr.Error()))
	}

	var err error

	// Handle different status transitions
	switch newStatus {
	case string(loan.EventApprove):
		err = h.loanService.ApproveLoan(loanEntity, req.ApprovalEmployeeID, req.DocumentID)
	// TODO: Complete the statuses
	// case "reject":
	// 	err = h.loanService.RejectLoan(loan, req.UpdatedBy)
	// case "disburse":
	// 	err = h.loanService.DisburseLoan(loan, req.UpdatedBy)
	// case "complete":
	// 	err = h.loanService.CompleteLoan(loan, req.UpdatedBy)
	default:
		return c.JSON(http.StatusBadRequest, response.Error("Unsupported status transition"))
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success("Loan status updated successfully"))
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
