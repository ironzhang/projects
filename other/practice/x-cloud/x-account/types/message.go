package types

type RegisterRequest struct {
	Type     string
	Name     string
	Password string
}

type LoginRequest struct {
	Type     string
	Name     string
	Password string
}
