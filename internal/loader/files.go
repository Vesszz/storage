package loader

import (
	"fmt"
	"io"
	"log/slog"
	"main/internal/models"
	"mime/multipart"
	"os"
)

type FileList struct {
	Filenames []string `json:"filenames"`
}

type PageData struct {
	Title      string
	FilesNames FileList
}

func (l *Loader) createFile(filename string, file multipart.File) error {
	//ext := filepath.Ext(filename)
	//photosExt := [7]string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff"}
	//fileDirector := ""
	out, err := os.Create(l.fsCfg.Path + filename)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return fmt.Errorf("saving file: %w", err)
	}
	return nil
}

func (l *Loader) Get(path string) (multipart.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("no file found: %w", err)
	}
	return file, nil
}

func (l *Loader) GetAll() (*FileList, error) {
	files, err := os.ReadDir(l.fsCfg.Path)
	if err != nil {
		return nil, fmt.Errorf("reading uploads dir: %w", err)
	}
	filenames := make([]string, len(files)-2)
	i := 0
	for _, file := range files {
		if file.Name() == "photos" || file.Name() == "videos" {
			continue
		}
		filenames[i] = file.Name()
		i++
	}
	response := &FileList{Filenames: filenames}
	return response, nil
}

// todo save with ok name not user's name
func (l *Loader) SaveFile(fileModel models.File, file multipart.File) (*models.File, error) {
	err := l.createFile(fileModel.Name, file)
	if err != nil {
		slog.Error("create file: %w", err)
		return nil, fmt.Errorf("create file: %w", err)
	}
	err = l.db.QueryRow("INSERT INTO files(user_id, time_created, name, times_viewed) VALUES ($1, $2, $3, $4) RETURNING id",
		fileModel.UserID, fileModel.TimeCreated, fileModel.Name, fileModel.TimesViewed).Scan(&fileModel.ID)
	if err != nil {
		slog.Error("insert file: %w", err)
		return nil, fmt.Errorf("insert file: %w", err)
	}
	return &fileModel, nil
}

func (l *Loader) IndexInfo() (PageData, error) {
	fileNames, err := l.GetAll()
	if err != nil {
		return PageData{}, fmt.Errorf("getting all files: %w", err)
	}
	data := PageData{
		Title:      "super title",
		FilesNames: *fileNames,
	}
	return data, nil
}
