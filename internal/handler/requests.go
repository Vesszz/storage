package handler

type SignUpRequest struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type SingUpResponse struct {
	AccessToken string `json:"access_token"`
}

type SignInRequest struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}
