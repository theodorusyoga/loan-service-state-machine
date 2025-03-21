package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/request"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/response"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

type LoanHandler struct {
	loanService *loan.LoanService
}

func NewLoanHandler(loanService *loan.LoanService) *LoanHandler {
	return &LoanHandler{
		loanService: loanService,
	}
}

func (h *LoanHandler) CreateLoan(c echo.Context) error {
	var req request.CreateLoanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request"))
	}

	loan, err := h.loanService.CreateLoan(c.Request().Context(), req.Amount)
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
