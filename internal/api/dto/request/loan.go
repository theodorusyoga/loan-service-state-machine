package request

type CreateLoanRequest struct {
	BorrowerID  string  `json:"borrowerId" validate:"required,uuid"`
	Amount      float64 `json:"amount" validate:"required,gt=0"`
	Rate        float64 `json:"rate" validate:"required,gt=0,lt=100"`
	ROI         float64 `json:"roi" validate:"required,gt=0,lt=100,roiLessThanRate"`
	Description string  `json:"description"`
}

type ApprovalRequest struct {
	ApprovedBy string `json:"approvedBy" validate:"required"`
}
