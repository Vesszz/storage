package loader

import (
	"database/sql"
	"fmt"
	"io"
	"main/internal/config"
	"main/internal/models"
	"mime/multipart"
	"os"

	_ "github.com/lib/pq"
)

type Loader struct {
	dbCfg *config.DatabaseConfig
	flCfg *config.FileStorageConfig
	db    *sql.DB
}

func New(dbCfg *config.DatabaseConfig, fsCfg *config.FileStorageConfig) (*Loader, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%d sslmode=%s", dbCfg.User, dbCfg.Password, dbCfg.DBName, dbCfg.Host, dbCfg.Port, dbCfg.SSLMode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("opening db connection: %w", err)
	}
	return &Loader{
		dbCfg: dbCfg,
		flCfg: fsCfg,
		db:    db,
	}, nil
}

type FileList struct {
	Filenames []string `json:"filenames"`
}

type PageData struct {
	Title      string
	FilesNames FileList
}

func (l *Loader) createFile(filename string) (*os.File, error) {
	//ext := filepath.Ext(filename)
	//photosExt := [7]string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".webp", ".tiff"}
	//fileDirector := ""
	return os.Create(l.flCfg.Path + filename)
}

func (l *Loader) Get(path string) (multipart.File, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("no file found: %w", err)
	}
	return file, nil
}

func (l *Loader) GetAll() (*FileList, error) {
	files, err := os.ReadDir(l.flCfg.Path)
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
	out, err := l.createFile(fileName)
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

func (l *Loader) CreateUser(user models.User) error {
	_, err := l.db.Exec("INSERT INTO users(name, password, salt) VALUES ($1, $2, $3)", user.Name, user.Password, user.Salt)
	if err != nil {
		return fmt.Errorf("inserting user: %w", err)
	}
	return nil
}
