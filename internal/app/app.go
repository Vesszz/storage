package app

import (
	"fmt"
	"main/internal/config"
	"main/internal/handler"
	"main/internal/logger"
	"main/internal/logic"
	"net/http"
)

func Run(configPath string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	loggerSync, err := logger.ReplaceZap(cfg.Logger)
	if err != nil {
		return fmt.Errorf("logger initialisation: %w", err)
	}
	defer loggerSync()
	l, err := logic.New(cfg)
	if err != nil {
		return fmt.Errorf("logic initialisation: %w", err)
	}
	h := handler.New(l)
	http.HandleFunc("/upload", h.Upload)
	//http.HandleFunc("/download/", h.Download)
	//http.HandleFunc("/list", h.FileList)
	http.HandleFunc("/signup", h.SignUp)
	http.HandleFunc("/signin", h.SignIn)
	http.HandleFunc("/refresh", h.RefreshTokens)
	//http.HandleFunc("/", h.Index)
	err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.Server.Port), nil)
	if err != nil {
		return fmt.Errorf("listening and serving: %w", err)
	}
	return nil
}
