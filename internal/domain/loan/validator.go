package loan

import (
	"errors"
	"fmt"
)

type DefaultStatusValidator struct{}

func NewDefaultStatusValidator() *DefaultStatusValidator {
	return &DefaultStatusValidator{}
}

// Check if the transition is valid
func (v *DefaultStatusValidator) isValidTransition(from, to Status) bool {
	validTransitions := map[Status][]Status{
		StatusProposed: {StatusApproved, StatusRejected},
		StatusApproved: {StatusInvested},
		StatusInvested: {StatusDisbursed},

		// Terminated
		StatusDisbursed: {},
		StatusRejected:  {},
	}

	allowedNext, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, status := range allowedNext {
		if status == to {
			return true
		}
	}

	return false
}

// Validate if loan can be approved
func (v *DefaultStatusValidator) validateApproval(loan *Loan) error {
	if loan.ApprovedBy == nil {
		return errors.New("approvedBy is required")
	}

	return nil
}

// Validate checks if a status transition is valid
func (v *DefaultStatusValidator) Validate(loan *Loan, from, to Status) error {
	// First check if the transition follows the state machine rules
	if !v.isValidTransition(from, to) {
		return fmt.Errorf("invalid transition from %s to %s", from, to)
	}

	// Then check business rules specific to each transition
	switch {
	case from == StatusProposed && to == StatusApproved:
		return v.validateApproval(loan)
	}

	return nil
}
