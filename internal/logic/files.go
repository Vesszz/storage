package logic

import (
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"main/internal/fileloader"
	"main/internal/models"
	"mime/multipart"
	"net/http"
	"time"
)

func (l *Logic) Index() (fileloader.PageData, error) {
	indexInfo, err := l.fileLoader.IndexInfo()
	if err != nil {
		return fileloader.PageData{}, fmt.Errorf("getting page info: %w", err)
	}
	return indexInfo, nil
}

func (l *Logic) Upload(userFileName string, fileName string, fileDescription string, file multipart.File, atc *AccessTokenClaims) error {
	key := uuid.New()
	_, err := l.fileLoader.SaveFile(models.File{
		UserID:      atc.ID,
		Key:         key,
		Path:        key.String() + fileName,
		TimeCreated: time.Now(),
		Name:        userFileName,
		Description: fileDescription,
		TimesViewed: 0,
	}, file)
	if err != nil {
		slog.Error("uploading file: %w", err)
		return fmt.Errorf("uploading file: %w", err)
	}
	return nil
}

// todo
func (l *Logic) Download(w http.ResponseWriter, fileName string) error {
	//file, err := l.loader.Get("./uploads/" + fileName)
	//if err != nil {
	//	return fmt.Errorf("getting file: %w", err)
	//}
	//w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
	//w.Header().Set("Content-Type", "application/octet-stream")
	//_, err = io.Copy(w, file)
	//if err != nil {
	//	return fmt.Errorf("downloading: %w", err)
	//}
	//defer file.Close()
	return nil
}

func (l *Logic) FileList() (*fileloader.FileList, error) {
	files, err := l.fileLoader.GetAll()
	if err != nil {
		return nil, fmt.Errorf("getting all files %w", err)
	}
	return files, nil
}
