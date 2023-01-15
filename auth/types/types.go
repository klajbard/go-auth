package types

type Status struct {
	Message string // `json:"message"`
}

type User struct {
	Username   string // `json:"username"`
	Password   string // `json:"password"`
	RememberMe bool   // `json:"rememberMe"`
}
