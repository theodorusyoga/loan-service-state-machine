package loanlender

import (
	"context"
	"errors"
)

var (
	ErrLoanLenderNotFound = errors.New("loan-lender relationship not found")
	ErrInvalidAmount      = errors.New("invalid investment amount")
)

type Service interface {
	Get(ctx context.Context, id string) (*LoanLender, error)
	GetByLoan(ctx context.Context, loanID string) ([]*LoanLender, error)
	Create(ctx context.Context, loanID, lenderID string, amount float64) (*LoanLender, error)
	List(ctx context.Context, filter LoanLenderFilter) ([]*LoanLender, error)
	TotalInvestmentForLoan(ctx context.Context, loanID string) (float64, error)
}

type LoanLenderService struct {
	repository Repository
}

func NewLoanLenderService(r Repository) *LoanLenderService {
	return &LoanLenderService{
		repository: r,
	}
}

func (s *LoanLenderService) Get(ctx context.Context, id string) (*LoanLender, error) {
	loanLender, err := s.repository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if loanLender == nil {
		return nil, ErrLoanLenderNotFound
	}
	return loanLender, nil
}

func (s *LoanLenderService) GetByLoan(ctx context.Context, loanID string) ([]*LoanLender, error) {
	return s.repository.GetByLoanID(ctx, loanID)
}

func (s *LoanLenderService) Create(ctx context.Context, loanID, lenderID string, amount float64) (*LoanLender, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	loanLender := NewLoanLender(loanID, lenderID, amount)
	if err := s.repository.Create(ctx, loanLender); err != nil {
		return nil, err
	}
	return loanLender, nil
}

func (s *LoanLenderService) List(ctx context.Context, filter LoanLenderFilter) ([]*LoanLender, error) {
	return s.repository.List(ctx, filter)
}

func (s *LoanLenderService) TotalInvestmentForLoan(ctx context.Context, loanID string) (float64, error) {
	investments, err := s.repository.GetByLoanID(ctx, loanID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, investment := range investments {
		total += investment.Amount
	}
	return total, nil
}
