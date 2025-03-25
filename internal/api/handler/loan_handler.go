package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/request"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/response"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/lender"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
)

type LoanHandler struct {
	loanService   *loan.LoanService
	lenderService *lender.LenderService
	validate      *validator.Validate
}

func NewLoanHandler(loanService *loan.LoanService, lenderService *lender.LenderService, validate *validator.Validate) *LoanHandler {
	return &LoanHandler{
		loanService:   loanService,
		lenderService: lenderService,
		validate:      validate,
	}
}

// ListLoans godoc
// @Summary List all loans
// @Description Get a list of all loans with optional filtering
// @Tags loans
// @Accept json
// @Produce json
// @Param max_amount query number false "Maximum loan amount"
// @Param min_amount query number false "Minimum loan amount"
// @Success 200 {object} domain.PaginatedResponse
// @Failure 500 {object} response.APIResponse
// @Router /loans [get]
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

// CreateLoan godoc
// @Summary Create a new loan
// @Description Create a new loan with the provided details
// @Tags loans
// @Accept json
// @Produce json
// @Param loan body request.CreateLoanRequest true "Loan information"
// @Success 201 {object} response.APIResponse
// @Failure 400 {object} response.APIResponse
// @Router /loans [post]
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

// UpdateLoanStatus godoc
// @Summary Update loan status
// @Description Update a loan's status based on the provided status transition
// @Description - For approve: { "success": true, "message": "Loan status updated successfully" }
// @Description - For partial invest: { "success": true, "data": { "remaining_amount": 150000, "invested_amount": 50000, "agreement_document": null }, "message": "loan invested successfully" }
// @Description - For full invest: { "success": true, "data": { "remaining_amount": 0, "invested_amount": 200000, "agreement_document": "agreement_file.pdf" }, "message": "loan status updated to invested" }
// @Description - For disburse: { "success": true, "data": { "field_officer_id": "emp-789", "agreement_file_name": "agreement.pdf" }, "message": "loan disbursed successfully" }
// @Tags loans
// @Accept json
// @Produce json
// @Param id path string true "Loan ID"
// @Param status path string true "New status"
// @Param request body object true "Status update information"
// @Param approveRequest body swagger.ApproveSchema false "Approve request (when status=approve)"
// @Param investRequest body swagger.InvestSchema false "Invest request (when status=invest)"
// @Param disburseRequest body swagger.DisburseSchema false "Disburse request (when status=disburse)"
// @Success 200 {object} response.APIResponse "Successful status update with varying response structure based on status"
// @Failure 400 {object} response.APIResponse "Invalid request or status transition"
// @Router /loans/{id}/{status} [patch]
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
		err = h.loanService.ApproveLoan(loanEntity, req.ApprovalEmployeeID, req.FileName)
	// TODO: Complete the statuses
	case string(loan.EventInvest):
		// check lender ID exists
		if req.LenderID == "" {
			return c.JSON(http.StatusBadRequest, response.Error("lender ID is required"))
		}
		// get lender
		lender, lenderErr := h.lenderService.GetByID(c.Request().Context(), req.LenderID)
		if lenderErr != nil {
			return c.JSON(http.StatusBadRequest, response.Error(lenderErr.Error()))
		}
		result, err := h.loanService.InvestLoan(loanEntity, lender, req.InvestAmount)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		}
		if result.RemainingAmount > 0 {
			return c.JSON(http.StatusOK, response.Success(result, "loan invested successfully"))
		} else {
			return c.JSON(http.StatusOK, response.Success(result, "loan status updated to invested"))
		}
	case string(loan.EventDisburse):
		result, err := h.loanService.DisburseLoan(loanEntity, req.FieldOfficerID, req.AgreementFileName)
		if err != nil {
			return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
		}
		return c.JSON(http.StatusOK, response.Success(result, "loan disbursed successfully"))

	default:
		return c.JSON(http.StatusBadRequest, response.Error("unsupported status transition"))
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, response.Success(nil, "Loan status updated successfully"))
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
