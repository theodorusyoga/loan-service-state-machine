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

	return c.JSON(http.StatusOK, response.Success(employees, "Employees retrieved successfully"))
}
