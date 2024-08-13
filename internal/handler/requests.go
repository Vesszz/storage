package handler

import (
	"mime/multipart"
)

type SignUpRequest struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type SignUpResponse struct {
	AccessToken string `json:"Authorization"`
}

type SignInRequest struct {
	Name        string `json:"name"`
	Password    string `json:"password"`
	Fingerprint string `json:"fingerprint"`
}

type UploadRequestForm struct {
	File        multipart.File
	FileName    string
	Description string
}

type UploadRequest struct {
	Description string `json:"description"`
}

type SignInResponse struct {
	AccessToken string `json:"Authorization"`
}

type RefreshRequest struct {
	Fingerprint string `json:"fingerprint"`
}

type RefreshResponse struct {
	AccessToken string `json:"Authorization"`
}
