package dto

type SignUp struct {
	FullName string `form:"fullName"`
	Email    string `form:"email"`
	Password string `form:"password"`
}
