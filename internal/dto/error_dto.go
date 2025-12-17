package dto

type ErrorResponse struct {
	Message string `json:"message" example:"Email is required field"`
}
