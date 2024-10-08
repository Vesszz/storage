package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"html/template"
	"main/internal/errs"
	"main/internal/logic"
	"net/http"
	"strings"
)

type Handler struct {
	logic *logic.Logic
}

func New(l *logic.Logic) *Handler {
	return &Handler{
		logic: l,
	}
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	data, err := h.logic.Index()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	file, fileHeader, err := r.FormFile("file")
	userFileName := r.Form.Get("name")
	userFileDescription := r.Form.Get("description")

	if err != nil {
		//http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		http.Error(w, "No file sent", http.StatusBadRequest)
		return
	}

	atc, err := h.decodeHeaderIntoAccessTokenClaims(r)
	if err != nil {
		http.Error(w, "Failed read access token", http.StatusBadRequest)
		return
	}

	err = h.logic.Upload(userFileName, fileHeader.Filename, userFileDescription, file, atc)

	if err != nil {
		//http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		zap.S().Errorf("failed to upload file: %v", err)
		http.Error(w, "Failed to upload file", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, r.Header.Get("Referer"), http.StatusOK)
}

func (h *Handler) Download(w http.ResponseWriter, r *http.Request) {
	fileName := r.URL.Path[10:] // 10 = /download/
	err := h.logic.Download(w, fileName)
	if err != nil {
		http.Error(w, "Sending error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) FileList(w http.ResponseWriter, r *http.Request) {
	fileNames, err := h.logic.FileList()
	if err != nil {
		http.Error(w, "Getting file uploaded", http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(fileNames)
	if err != nil {
		http.Error(w, "Encoding json", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var singUpReq SignUpRequest
	err := json.NewDecoder(r.Body).Decode(&singUpReq)
	if err != nil {
		zap.S().Errorf("decoding json: %v", err)
		http.Error(w, "Decoding json", http.StatusBadRequest)
		return
	}

	tokens, err := h.logic.Register(singUpReq.Name, singUpReq.Password, singUpReq.Fingerprint)
	if err != nil {
		if errors.Is(err, errs.UserAlreadyExists) {
			http.Error(w, "This name is taken", http.StatusBadRequest)
			return
		}
		zap.S().Errorf("managing reg: %v", err)
		http.Error(w, "Managing registration", http.StatusInternalServerError)
		return
	}

	setRefreshTokenToCookie(w, tokens.RefreshToken.String())
	response := SignUpResponse{AccessToken: tokens.AccessToken}
	err = encodeResponse(w, response)
	if err != nil {
		zap.S().Errorf("encoding response: %v", err)
		http.Error(w, "Encoding response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var singInReq SignInRequest
	err := json.NewDecoder(r.Body).Decode(&singInReq)
	if err != nil {
		http.Error(w, "Decoding json", http.StatusBadRequest)
		return
	}

	tokens, err := h.logic.Login(singInReq.Name, singInReq.Password, singInReq.Fingerprint)
	if err != nil {
		if errors.Is(err, errs.WrongPassword) {
			http.Error(w, "Wrong password", http.StatusBadRequest)
			return
		}
		zap.S().Errorf("login: %v", err)
		http.Error(w, "Signing in", http.StatusInternalServerError)
		return
	}

	setRefreshTokenToCookie(w, tokens.RefreshToken.String())
	response := SignInResponse{AccessToken: tokens.AccessToken}
	err = encodeResponse(w, response)
	if err != nil {
		zap.S().Errorf("encoding response: %v", err)
		http.Error(w, "Encoding response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	var refreshReq RefreshRequest
	err := json.NewDecoder(r.Body).Decode(&refreshReq)
	if err != nil {
		http.Error(w, "Decoding json", http.StatusBadRequest)
		return
	}
	refreshToken, err := getRefreshTokenFromRequest(r)
	if err != nil {
		if errors.Is(err, errs.RefreshTokenNotFound) {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}
		zap.S().Errorf("Get Refresh token from request: %v", err)
		http.Error(w, "Get Refresh token from request", http.StatusInternalServerError)
		return
	}
	tokens, err := h.logic.Refresh(refreshToken, refreshReq.Fingerprint)
	if err != nil {
		//todo
		http.Error(w, "Refreshing tokens", http.StatusInternalServerError)
		return
	}
	setRefreshTokenToCookie(w, tokens.RefreshToken.String())
	response := RefreshResponse{AccessToken: tokens.AccessToken}
	err = encodeResponse(w, response)
	if err != nil {
		zap.S().Errorf("encoding response: %v", err)
		http.Error(w, "Encoding response", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func getRefreshTokenFromRequest(r *http.Request) (string, error) {
	cookies := r.Cookies()

	var refreshToken string
	for _, cookie := range cookies {
		if cookie.Name == "refresh_token" {
			refreshToken = cookie.Value
			break
		}
	}

	if refreshToken == "" {
		return "", errs.RefreshTokenNotFound
	}
	return refreshToken, nil
}

func setRefreshTokenToCookie(w http.ResponseWriter, refreshToken string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Path:     "/",
		//Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})
}

func encodeResponse(w http.ResponseWriter, response any) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		return fmt.Errorf("encoding response: %w", err)
	}
	return nil
}

func (h *Handler) decodeHeaderIntoAccessTokenClaims(r *http.Request) (*logic.AccessTokenClaims, error) {
	accessToken := r.Header.Get("Authorization")
	if accessToken == "" {
		return nil, fmt.Errorf("access token not found")
	}
	ats := strings.Split(accessToken, " ")
	if len(ats) != 2 {
		return nil, fmt.Errorf("authorization failed, bad access token")
	}
	if ats[0] != "Bearer" {
		return nil, fmt.Errorf("authorization failed, bad access token")
	}
	atc, err := h.logic.AccessTokenClaimsFromAccessToken(ats[1])
	if err != nil {
		zap.S().Errorf("turn access token (string) to claims: %v", err)
		return nil, fmt.Errorf("turn access token (string) to claims: %w", err)
	}
	return atc, nil
}
