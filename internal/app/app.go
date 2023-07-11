package app

import (
	"time"

	"github.com/julinserg/OtusAlgorithmHomeProject/internal/storage"
)

type App struct {
	logger  Logger
	storage Storage
}

type DocumentSrc struct {
	Index int
	Url   string
}

type DocumentSearch struct {
	Index   int
	Url     string
	Context string
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

func (a *App) AddNewDocument(url string) ([]DocumentSrc, error) {
	documents := make([]DocumentSrc, 0)
	documents = append(documents, DocumentSrc{1, "https://www.w3schools.com/howto/howto_css_searchbar.asp"})
	documents = append(documents, DocumentSrc{2, "https://ru.wikipedia.org/wiki/Yahoo!_Search"})
	return documents, nil
}

func (a *App) Search(str string) ([]DocumentSearch, error) {
	documents := make([]DocumentSearch, 0)
	documents = append(documents, DocumentSearch{
		1,
		"https://stackoverflow.com/questions/9523927/how-to-stop-table-from-resizing-when-contents-grow",
		"I have a table, the cells of which are filled with picture"})
	documents = append(documents, DocumentSearch{
		2,
		"https://stackoverflow.com/questions/21019302/html-button-layout-positioning",
		"Even i didn't get what exactly you want. but for an image sourrounded by buttons try this code"})
	return documents, nil
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}
