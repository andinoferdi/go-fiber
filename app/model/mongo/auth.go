package model

type LoginRequest struct {
	Email    string `bson:"email" json:"email" validate:"required,email"`
	Password string `bson:"password" json:"password" validate:"required"`
}

type LoginResponse struct {
	Alumni Alumni `bson:"alumni" json:"alumni"`
	Token  string `bson:"token" json:"token"`
}

type GetProfileResponse struct {
	Success bool   `bson:"success" json:"success"`
	Message string `bson:"message" json:"message"`
	Data    struct {
		AlumniID string `bson:"alumni_id" json:"alumni_id"`
		Email    string `bson:"email" json:"email"`
		Role     string `bson:"role" json:"role"`
	} `bson:"data" json:"data"`
}
