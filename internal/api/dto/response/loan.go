package response

import "time"

type LoanResponse struct {
	ID                string                `json:"id"`
	Amount            float64               `json:"amount"`
	Status            string                `json:"status"`
	CreatedAt         time.Time             `json:"createdAt"`
	UpdatedAt         time.Time             `json:"updatedAt"`
	ApprovalDate      *time.Time            `json:"approvalDate,omitempty"`
	ApprovedBy        string                `json:"approvedBy,omitempty"`
	InvestmentDate    *time.Time            `json:"investmentDate,omitempty"`
	DisbursementDate  *time.Time            `json:"disbursementDate,omitempty"`
	StatusTransitions []StatusTransitionDTO `json:"statusTransitions"`
}

type StatusTransitionDTO struct {
	From        string    `json:"from"`
	To          string    `json:"to"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	PerformedBy string    `json:"performedBy"`
}

type LoanLenderResponse struct {
	RemainingAmount   float64 `json:"remaining_amount"`
	InvestedAmount    float64 `json:"invested_amount"`
	AgreementDocument *string `json:"agreement_document"`
}

type DisbursementResponse struct {
	DisbursementDate  time.Time `json:"disbursement_date"`
	DisbursedBy       string    `json:"disbursed_by"`
	AgreementDocument *string   `json:"agreement_document"`
	BorrowerRepayment float64   `json:"borrower_repayment"`
	InvestorROI       float64   `json:"investor_roi"`
}
