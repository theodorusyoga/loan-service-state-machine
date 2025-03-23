package request

type CreateBorrowerRequest struct {
	FullName    string `json:"fullName" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phoneNumber" validate:"required"`
	IDNumber    string `json:"idNumber" validate:"required"`
}
