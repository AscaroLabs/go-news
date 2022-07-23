package storage

type UserDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type TokenDTO struct {
	UserId string
	Role   string
}

type SignInDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
