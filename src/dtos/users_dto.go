package dtos

type UserCreateDTO struct {
	Name      string  `json:"name" validate:"required"`
	Username  string  `json:"username" validate:"required"`
	Email     string  `json:"email" validate:"required,email"`
	Password  string  `json:"password" validate:"required,min=6"`
	CreatedBy *string `json:"created_by" validate:"required"`
	UpdatedBy *string `json:"updated_by"`
}
