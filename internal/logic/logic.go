package logic

import (
	"fmt"
	"io"
	"main/internal/loader"
	"mime/multipart"
	"net/http"
)

type Logic struct {
	loader *loader.Loader
}

func New(l *loader.Loader) *Logic {
	return &Logic{
		loader: l,
	}
}

func (l *Logic) Index() (loader.PageData, error) {
	indexInfo, err := l.loader.IndexInfo()
	if err != nil {
		return loader.PageData{}, fmt.Errorf("getting page info: %w", err)
	}
	return indexInfo, nil
}

func (l *Logic) Upload(fileName string, file multipart.File) error {
	err := l.loader.Load(fileName, file)
	if err != nil {
		return fmt.Errorf("loading file: %w", err)
	}
	return nil
}

func (l *Logic) Download(w http.ResponseWriter, fileName string) error {
	file, err := l.loader.Get("./uploads/" + fileName)
	if err != nil {
		return fmt.Errorf("getting file: %w", err)
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, file)
	if err != nil {
		return fmt.Errorf("downloading: %w", err)
	}
	defer file.Close()
	return nil
}

func (l *Logic) FileList() (*loader.FileList, error) {
	files, err := l.loader.GetAll()
	if err != nil {
		return nil, fmt.Errorf("getting all files %w", err)
	}
	return files, nil
}
