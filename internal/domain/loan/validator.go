package loan

import (
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

// Validate checks if a status transition is valid
func (v *DefaultStatusValidator) Validate(loan *Loan, from, to Status) error {
	if !v.isValidTransition(from, to) {
		return fmt.Errorf("cannot change status from %s to %s", from, to)
	}

	return nil
}
