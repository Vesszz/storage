package loader

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

type Loader struct {
}

func New() *Loader {
	return &Loader{}
}

type FileList struct {
	Filenames []string `json:"filenames"`
}

type PageData struct {
	Title      string
	FilesNames FileList
}

type Config struct {
	Database struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbname"`
		SSLMode  string `json:"sslmode"`
	} `json:"database"`
}

func createFile(filename string) (*os.File, error) {
	//ext := filepath.Ext(filename)
	//photosExt := [7]string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff"}
	//fileDirector := ""
	return os.Create("./uploads/" + filename)
}

func (l *Loader) Get(path string) (multipart.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("no file found: %w", err)
	}
	return file, nil
}

func (l *Loader) GetAll() (*FileList, error) {
	files, err := os.ReadDir("./uploads")
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

func (l *Loader) Load(fileName string, file multipart.File) error {
	out, err := createFile(fileName)
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

func (l *Loader) DBstart() error {
	file, err := os.ReadFile("./configs/config.json")
	if err != nil {
		fmt.Println(err)
		return err
	}
	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err)
		return err
	}
	dbcfg := config.Database
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", dbcfg.User, dbcfg.Password, dbcfg.DBName, dbcfg.Host, dbcfg.Port, dbcfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	return nil
}
