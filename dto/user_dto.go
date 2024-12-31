package dto

type RegisterationDto struct {
	FirstName string `json:"first_name" binding:"required, min=3"`
	LastName  string `json:"last_name" binding:"required, min=3"`
	Username  string `json:"username" binding:"required, min=3"`
	Email     string `json:"email" binding:"required, email"` //email is gin builtin validation
	Password  string `json:"password" binding:"required, password"`
}

type LoginByUsernameDto struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required, password"`
}

type RegisterLoginByMobileDto struct {
	Mobile string `json:"mobile" binding:"required, mobile, min=11, max=11"`
	Otp    string `json:"otp" binding:"min=6, max=6"`
}
