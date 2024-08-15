package fileloader

import (
	"fmt"
	"go.uber.org/zap"
	"io"
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

func (l *FileLoader) createFile(filename string, file multipart.File) error {
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

func (l *FileLoader) GetByPath(path string) (multipart.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("no file found: %w", err)
	}
	return file, nil
}

func (l *FileLoader) GetAll() (*FileList, error) {
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

func (l *FileLoader) SaveFile(fileModel models.File, file multipart.File) (*models.File, error) {
	err := l.createFile(fileModel.Path, file)
	if err != nil {
		zap.S().Errorf("create file: %v", err)
		return nil, fmt.Errorf("create file: %w", err)
	}
	return &fileModel, nil
}

func (l *FileLoader) IndexInfo() (PageData, error) {
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
