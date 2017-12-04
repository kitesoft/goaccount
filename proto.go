package main

type LoginRequest struct {
	UserName string `json:"username" valid:"required"`
	Password string `json:"password" valid:"required,length(6|100)"`
}

type LoginResponse struct {
	UserId   uint32 `json:"user_id" valid:"required"`
	UserName string `json:"username"`
	RealName string `json:"real_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Token    string `json:"token"`
}

type UpdatePasswordRequest struct {
	UserName    string `json:"username",valid:"required"`
	NewPassword string `json:"new_password" valid:"required,length(6|100)"`
	OldPassword string `json:"old_password" valid:"required,length(6|100)"`
}

type ResetPasswordRequest struct {
	UserID []int `json:"user_id"`
}
