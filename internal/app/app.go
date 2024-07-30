package app

import (
	"fmt"
	"main/internal/config"
	"main/internal/handler"
	"main/internal/logic"
	"net/http"
)

func Run() error {
	cfg, err := config.Load("./configs/config.yaml")
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	l, err := logic.New(cfg)
	if err != nil {
		return fmt.Errorf("logic initialisation: %w", err)
	}
	h := handler.New(l)
	http.HandleFunc("/", h.Index)
	http.HandleFunc("/upload", h.Upload)
	http.HandleFunc("/download/", h.Download)
	http.HandleFunc("/list", h.FileList)
	http.HandleFunc("/singup", h.Registration)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil)
	if err != nil {
		return fmt.Errorf("listening and serving: %w", err)
	}
	return nil
}
