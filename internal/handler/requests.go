package handler

type SignUpRequest struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type SignUpResponse struct {
	AccessToken string `json:"access_token"`
}

type SignInRequest struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type SignInResponse struct {
	AccessToken string `json:"access_token"`
}

type RefreshRequest struct {
	Fingerprint string `json:"fingerprint"`
}

type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
