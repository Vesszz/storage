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

	err = h.logic.Upload(header.Filename, file)
	if err != nil {
		//http.Redirect(w, r, r.Header.Get("Referer"), http.StatusFound)
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

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken.String(),
		HttpOnly: true,
		Path:     "/",
		//Secure:   true,
		SameSite: http.SameSiteNoneMode,
	})

	response := SingUpResponse{AccessToken: tokens.AccessToken}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		slog.Error(fmt.Errorf("encoding json response: %w", err).Error())
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
	//TODO get token
	//tokens, err := h.logic.Login(singInReq.Name, singInReq.Password)

	return
}
