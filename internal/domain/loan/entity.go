package loan

import (
	"time"
)

type Status string

const (
	StatusProposed  Status = "proposed"
	StatusApproved  Status = "approved"
	StatusInvested  Status = "invested"
	StatusDisbursed Status = "disbursed"
	StatusRejected  Status = "rejected"
)

// To mark the history of the status transition
type StatusTransition struct {
	From        Status    `json:"from"`
	To          Status    `json:"to"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	PerformedBy string    `json:"performed_by"`
}

type Loan struct {
	ID                  string             `json:"id"`
	BorrowerID          string             `json:"borrower_id"`
	Amount              float64            `json:"amount"`
	Rate                float64            `json:"rate"`
	ROI                 float64            `json:"roi"`
	Status              Status             `json:"status"`
	SurveyDocumentID    string             `json:"survey_document_id"`
	ApprovalDate        *time.Time         `json:"approval_date"`
	ApprovedBy          *string            `json:"approved_by"`
	InvestmentDate      *time.Time         `json:"investment_date"`
	DisbursementDate    *time.Time         `json:"disbursement_date"`
	DisbursedBy         *string            `json:"disbursed_by"`
	AgreementDocumentID *string            `json:"agreement_document_id"`
	StatusTransitions   []StatusTransition `json:"status_transitions"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
}

func NewLoan(id string, borrowerID string, amount float64, rate float64, roi float64) *Loan {
	now := time.Now()

	return &Loan{
		ID:         id,
		BorrowerID: borrowerID,
		Amount:     amount,
		Rate:       rate,
		ROI:        roi,
		Status:     StatusProposed,
		StatusTransitions: []StatusTransition{
			{
				From:        "",
				To:          StatusProposed,
				Date:        time.Now(),
				Description: "Loan created",
				PerformedBy: "system",
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}
}
