package model

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Alumni Alumni `json:"alumni"`
	Token  string `json:"token"`
}

type GetProfileResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    struct {
		AlumniID int    `json:"alumni_id"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"data"`
}