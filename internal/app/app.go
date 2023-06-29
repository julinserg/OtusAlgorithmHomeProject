package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/julinserg/OtusAlgorithmHomeProject/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type Logger interface {
	Error(msg string)
}

type Storage interface {
	Add(event storage.Event) error
	Update(event storage.Event) error
	Remove(id string) error
	GetEventsByDay(date time.Time) ([]storage.Event, error)
	GetEventsByWeek(dateBeginWeek time.Time) ([]storage.Event, error)
	GetEventsByMonth(dateBeginMonth time.Time) ([]storage.Event, error)
}

type InputParams struct {
	Width    int
	Height   int
	ImageURL string
}

var ErrFromRemoteServer = errors.New("error from remote server")

func (a *App) getImageFromRemoteServer(
	imageURL string,
	header http.Header,
) ([]byte, int, error) {
	return nil, http.StatusOK, nil
}

func (a *App) cropAndResizeImage(imageRaw []byte, width int, height int) ([]byte, error) {
	return nil, nil
}

func (a *App) saveImageOnDisk(image []byte, pathToFile string) error {
	return nil
}

func (a *App) readImageFromDisk(pathToFile string) ([]byte, error) {
	return nil, nil
}

func (a *App) GetImagePreview(params InputParams, header http.Header) ([]byte, int, bool, error) {
	return nil, 0, true, nil
}

func (a *App) ClearCache() {

}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}
