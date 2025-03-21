package request

type CreateLoanRequest struct {
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Description string  `json:"description"`
}

type ApprovalRequest struct {
	ApprovedBy string `json:"approvedBy" validate:"required"`
}
