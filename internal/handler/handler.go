package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
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

func (h *Handler) Registration(w http.ResponseWriter, r *http.Request) {
	var regReq RegistrationRequest
	err := json.NewDecoder(r.Body).Decode(&regReq)
	if err != nil {
		http.Error(w, "Decoding json", http.StatusBadRequest)
		return
	}
	err = h.logic.Register(regReq.Name, regReq.Password)
	if err != nil {
		http.Error(w, "Managing registration", http.StatusInternalServerError)
		fmt.Println(err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
