package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
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
	file, header, err := r.FormFile("file")

	if err != nil {
		//http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		http.Error(w, "No file sent", http.StatusBadRequest)
		return
	}

	atc, err := h.decodeHeaderIntoAccessTokenClaims(r)
	if err != nil {
		slog.Error("failed read access token: %w", err)
		http.Error(w, "Failed read access token", http.StatusBadRequest)
		return
	}
	err = h.logic.Upload(header.Filename, file, atc)
	if err != nil {
		//http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
		slog.Error("failed to upload file: %w", err)
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
		slog.Error(fmt.Errorf("decoding json: %w", err).Error())
		http.Error(w, "Decoding json", http.StatusBadRequest)
		return
	}

	tokens, err := h.logic.Register(singUpReq.Name, singUpReq.Password, singUpReq.Fingerprint)
	if err != nil {
		if errors.Is(err, errs.UserAlreadyExists) {
			http.Error(w, "This name is taken", http.StatusBadRequest)
			return
		}
		slog.Error(fmt.Errorf("managing reg: %w", err).Error())
		http.Error(w, "Managing registration", http.StatusInternalServerError)
		return
	}

	setRefreshTokenToCookie(w, tokens.RefreshToken.String())
	response := SignUpResponse{AccessToken: tokens.AccessToken}
	err = encodeResponse(w, response)
	if err != nil {
		slog.Error("encoding response: %w", err)
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
		slog.Error(fmt.Errorf("login: %w", err).Error())
		http.Error(w, "Signing in", http.StatusInternalServerError)
		return
	}

	setRefreshTokenToCookie(w, tokens.RefreshToken.String())
	response := SignInResponse{AccessToken: tokens.AccessToken}
	err = encodeResponse(w, response)
	if err != nil {
		slog.Error("encoding response: %w", err)
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
		slog.Error("Get Refresh token from request: %w", err)
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
		slog.Error("encoding response: %w", err)
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
	//todo anti [1]
	atc, err := h.logic.AccessTokenClaimsFromAccessToken(strings.Split(accessToken, " ")[1])
	if err != nil {
		slog.Error("turn access token (string) to claims: %w", err)
		return nil, fmt.Errorf("turn access token (string) to claims: %w", err)
	}
	return atc, nil
}
