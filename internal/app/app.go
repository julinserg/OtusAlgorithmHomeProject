package app

import (
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
	Add(document storage.DocumentSource) error
	GetAllDocumentSource() ([]storage.DocumentSource, error)
}

func documentSrcFromStorage(doc *storage.DocumentSource) DocumentSrc {
	docApp := DocumentSrc{
		Index: doc.Index, Url: doc.Url,
	}
	return docApp
}

func documentSrcToStorage(docApp *DocumentSrc) storage.DocumentSource {
	docStor := storage.DocumentSource{
		Index: docApp.Index, Url: docApp.Url,
	}
	return docStor
}

func (a *App) AddNewDocument(url string) ([]DocumentSrc, error) {
	documents := make([]DocumentSrc, 0)
	a.storage.Add(storage.DocumentSource{Url: url})
	listDoc, err := a.storage.GetAllDocumentSource()
	if err != nil {
		return nil, err
	}
	for _, d := range listDoc {
		documents = append(documents, DocumentSrc{Index: d.Index, Url: d.Url})
	}
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
