package callbacks

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"unsafe"

	"github.com/looplab/fsm"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/loan/callbacks"
	"github.com/theodorusyoga/loan-service-state-machine/internal/test/mocks"
)

// To assign private field `cancelFunc` in fsm.Event
func setCancelFunc(e *fsm.Event, fn func()) {
	val := reflect.ValueOf(e).Elem()
	field := val.FieldByName("cancelFunc")

	// Create a pointer to the field
	ptr := unsafe.Pointer(field.UnsafeAddr())

	// Create a reflect.Value of the function
	fnVal := reflect.ValueOf(fn)

	// Set the pointer to the function
	reflect.NewAt(field.Type(), ptr).Elem().Set(fnVal)
}

func TestBeforeApprove(t *testing.T) {
	t.Run("should pass when all validations succeed", func(t *testing.T) {
		// Setup
		mockEmployeeRepo := mocks.NewMockEmployeeRepository()

		provider := &callbacks.CallbackProvider{
			Validator:          loan.DefaultStatusValidator{},
			EmployeeRepository: mockEmployeeRepo,
		}

		loanObj := &loan.Loan{ID: "loan-123"}
		approvedBy := "employee-123"
		fileName := "document.pdf"

		mockEmployeeRepo.On("Get", mock.Anything, approvedBy).Return(struct{}{}, nil)

		// Create mock event
		mockEvent := &fsm.Event{
			Src:  "proposed",
			Dst:  "approved",
			Args: []interface{}{loanObj, approvedBy, fileName},
			FSM:  &fsm.FSM{},
		}

		// Execute
		provider.BeforeApproval(context.Background(), mockEvent)

		// Assert
		mockEmployeeRepo.AssertExpectations(t)
	})

	t.Run("should cancel when document is missing", func(t *testing.T) {
		mockEmployeeRepo := mocks.NewMockEmployeeRepository()

		provider := &callbacks.CallbackProvider{
			Validator:          loan.DefaultStatusValidator{},
			EmployeeRepository: mockEmployeeRepo,
		}

		loanObj := &loan.Loan{ID: "loan-123"}
		approvedBy := "employee-123"
		fileName := "" // Empty filename

		mockEmployeeRepo.On("Get", mock.Anything, approvedBy).Return(struct{}{}, nil)

		// Create event
		mockEvent := &fsm.Event{
			Src:  "proposed",
			Dst:  "approved",
			Args: []interface{}{loanObj, approvedBy, fileName},
			FSM:  &fsm.FSM{},
		}

		setCancelFunc(mockEvent, func() {})

		// Execute
		provider.BeforeApproval(context.Background(), mockEvent)

		// Assert
		assert.Equal(t, "document is required", mockEvent.Err.Error())
	})

	t.Run("should cancel when approver is missing", func(t *testing.T) {
		// Setup
		provider := &callbacks.CallbackProvider{}

		loanObj := &loan.Loan{ID: "loan-123"}
		approvedBy := "" // Empty approver
		fileName := "document.pdf"

		// Create event
		mockEvent := &fsm.Event{
			Args: []interface{}{loanObj, approvedBy, fileName},
		}

		setCancelFunc(mockEvent, func() {})

		// Execute
		provider.BeforeApproval(context.Background(), mockEvent)

		// Assert
		assert.Equal(t, "approved by is required", mockEvent.Err.Error())
	})

	t.Run("should cancel when transition validation fails", func(t *testing.T) {
		// Setup
		mockValidator := new(MockValidator)
		provider := &callbacks.CallbackProvider{
			Validator: *loan.NewDefaultStatusValidator(),
		}

		loanObj := &loan.Loan{ID: "loan-123"}
		approvedBy := "employee-123"
		fileName := "document.pdf"

		// Configure mocks
		validationError := errors.New("cannot change status from approved to proposed")
		mockValidator.On("Validate", loanObj, loan.Status("approved"), loan.Status("proposed")).Return(validationError)

		// Create event
		mockEvent := &fsm.Event{
			Src:  "approved",
			Dst:  "proposed",
			Args: []interface{}{loanObj, approvedBy, fileName},
		}

		setCancelFunc(mockEvent, func() {})

		// Execute
		provider.BeforeApproval(context.Background(), mockEvent)

		// Assert
		assert.EqualValues(t, validationError.Error(), mockEvent.Err.Error())
	})

	t.Run("should cancel when employee is not found", func(t *testing.T) {
		// Setup
		mockValidator := new(MockValidator)
		mockEmployeeRepo := mocks.NewMockEmployeeRepository()

		provider := &callbacks.CallbackProvider{
			Validator:          *loan.NewDefaultStatusValidator(),
			EmployeeRepository: mockEmployeeRepo,
		}

		loanObj := &loan.Loan{ID: "loan-123"}
		approvedBy := "employee-123"
		fileName := "document.pdf"

		// Configure mocks
		mockValidator.On("Validate", loanObj, loan.Status("source-state"), loan.Status("dest-state")).Return(nil)
		mockEmployeeRepo.On("Get", mock.Anything, approvedBy).Return(nil, errors.New("employee not found"))

		// Create event
		mockEvent := &fsm.Event{
			Src:  "proposed",
			Dst:  "approved",
			Args: []interface{}{loanObj, approvedBy, fileName},
		}

		setCancelFunc(mockEvent, func() {})

		// Execute
		provider.BeforeApproval(context.Background(), mockEvent)

		// Assert
		assert.Equal(t, "employee not found", mockEvent.Err.Error())
		mockEmployeeRepo.AssertExpectations(t)
	})
}
