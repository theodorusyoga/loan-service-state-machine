package request

type CreateDocumentRequest struct {
	LoanID   string `json:"loanId" validate:"required"`
	FileName string `json:"fileName" validate:"required"`
}
