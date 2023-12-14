package dto

type Login struct {
	Email string `form:"email"`
	Password string `form:"password"`
}
