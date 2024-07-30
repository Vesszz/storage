package handler

type RegistrationRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
