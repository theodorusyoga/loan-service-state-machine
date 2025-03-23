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
	From        Status
	To          Status
	Date        time.Time
	Description string
	PerformedBy string
}

type Loan struct {
	ID                  string
	BorrowerID          string
	Amount              float64
	Rate                float64
	ROI                 float64
	Status              Status
	SurveyDocumentID    string
	ApprovalDate        *time.Time
	ApprovedBy          *string
	InvestmentDate      *time.Time
	DisbursementDate    *time.Time
	DisbursedBy         *string
	AgreementDocumentID string
	StatusTransitions   []StatusTransition
	CreatedAt           time.Time
	UpdatedAt           time.Time
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
