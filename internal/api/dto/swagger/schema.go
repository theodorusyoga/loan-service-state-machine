package swagger

// This file contains Swagger schema definitions for the API documentation

// LoanStatusUpdateSchemas contains documentation for the various request bodies
// used with loan status updates based on the status parameter
type LoanStatusUpdateSchemas struct {
	// These fields are for documentation purposes only and are not used in code
}

// ApproveSchema defines the request structure for loan approval (status=approve)
type ApproveSchema struct {
	ApprovalEmployeeID string `json:"approval_employee_id" example:"emp-123"`
	ApprovalDate       string `json:"approval_date" example:"2025-03-25"`
	FileName           string `json:"file_name" example:"approval_document.pdf"`
}

// InvestSchema defines the request structure for loan investment (status=invest)
type InvestSchema struct {
	LenderID     string  `json:"lender_id" example:"lender-456"`
	InvestAmount float64 `json:"invest_amount" example:"5000.00"`
}

// DisburseSchema defines the request structure for loan disbursement (status=disburse)
type DisburseSchema struct {
	FieldOfficerID    string `json:"field_officer_id" example:"emp-789"`
	AgreementFileName string `json:"agreement_file_name" example:"loan_agreement.pdf"`
}

type SuccessResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Loan status updated successfully"`
}

// PartialInvestResponse represents a successful partial investment response
type PartialInvestResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"loan invested successfully"`
	Data    struct {
		RemainingAmount   float64 `json:"remaining_amount" example:"150000"`
		InvestedAmount    float64 `json:"invested_amount" example:"50000"`
		AgreementDocument string  `json:"agreement_document" example:"agreement-doc.pdf"`
	} `json:"data"`
}

// FullInvestResponse represents a response when a loan is fully invested
type FullInvestResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"loan status updated to invested"`
	Data    struct {
		RemainingAmount   float64 `json:"remaining_amount" example:"0"`
		InvestedAmount    float64 `json:"invested_amount" example:"200000"`
		AgreementDocument string  `json:"agreement_document" example:"agreement_effa1b23-489c-4638-97e6-7bed3e15006c.pdf"`
	} `json:"data"`
}

// DisburseResponse represents a successful disbursement response
type DisburseResponse struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"loan disbursed successfully"`
	Data    struct {
		FieldOfficerID    string `json:"field_officer_id" example:"emp-789"`
		AgreementFileName string `json:"agreement_file_name" example:"agreement_123456.pdf"`
		DisbursementDate  string `json:"disbursement_date,omitempty" example:"2025-03-25T15:30:45Z"`
	} `json:"data"`
}
