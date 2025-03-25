package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/request"
	"github.com/theodorusyoga/loan-service-state-machine/internal/api/dto/response"
	"github.com/theodorusyoga/loan-service-state-machine/internal/domain/employee"
)

type EmployeeHandler struct {
	employeeService *employee.EmployeeService
	validate        *validator.Validate
}

func NewEmployeeHandler(employeeService *employee.EmployeeService, validate *validator.Validate) *EmployeeHandler {
	return &EmployeeHandler{
		employeeService: employeeService,
		validate:        validate,
	}
}

// CreateEmployee godoc
// @Summary Create a new employee
// @Description Register a new employee in the system
// @Tags employees
// @Accept json
// @Produce json
// @Param employee body request.CreateEmployeeRequest true "Employee information"
// @Success 201 {object} response.APIResponse{data=employee.Employee} "Employee created successfully"
// @Failure 400 {object} response.APIResponse "Invalid request or validation error"
// @Router /employees [post]
func (h *EmployeeHandler) CreateEmployee(c echo.Context) error {
	var req request.CreateEmployeeRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, response.Error("Invalid request"))
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		// Format validation errors nicely
		validationErrors := err.(validator.ValidationErrors)
		errorsMsg := formatValidationErrors(validationErrors)
		return c.JSON(http.StatusBadRequest, response.Error(errorsMsg))
	}

	employee, err := h.employeeService.CreateEmployee(c.Request().Context(), req.FullName, req.Email, req.PhoneNumber, req.IDNumber)
	if err != nil {
		return c.JSON(http.StatusBadRequest, response.Error(err.Error()))
	}

	return c.JSON(http.StatusCreated, response.Success(employee, "Employee created successfully"))
}

// ListEmployees godoc
// @Summary List all employees
// @Description Get a list of all employees with optional filtering
// @Tags employees
// @Accept json
// @Produce json
// @Param full_name query string false "Filter by full name"
// @Param email query string false "Filter by email"
// @Param phone_number query string false "Filter by phone number"
// @Param id_number query string false "Filter by ID number"
// @Success 200 {object} domain.PaginatedResponse{data=[]employee.Employee} "List of employees"
// @Failure 500 {object} response.APIResponse "Internal server error"
// @Router /employees [get]
func (h *EmployeeHandler) ListEmployees(c echo.Context) error {
	// Extract query parameters for filtering
	fullName := c.QueryParam("full_name")
	email := c.QueryParam("email")
	phoneNumber := c.QueryParam("phone_number")
	idNumber := c.QueryParam("id_number")

	filter := employee.EmployeeFilter{
		FullName:    &fullName,
		Email:       &email,
		PhoneNumber: &phoneNumber,
		IDNumber:    &idNumber,
	}

	employees, err := h.employeeService.ListEmployees(c.Request().Context(), filter)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, response.Error(err.Error()))
	}

	return c.JSON(http.StatusOK, employees)
}
