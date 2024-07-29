package app

import (
	"fmt"
	"main/internal/handler"
	"main/internal/loader"
	"main/internal/logic"
	"net/http"
)

func Run() error {
	load := loader.New()
	log := logic.New(load)
	h := handler.New(log)
	http.HandleFunc("/", h.Index)
	http.HandleFunc("/upload", h.Upload)
	http.HandleFunc("/download/", h.Download)
	http.HandleFunc("/list", h.FileList)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return fmt.Errorf("listening and serving: %w", err)
	}
	return nil
}
