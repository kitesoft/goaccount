package model

type UserState struct {
	UserId   int32  `json:"user_id"`
	UserName string `json:"user_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	RealName string `json:"real_name"`
	Country  string `json:"country"`
}
