package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/julinserg/OtusAlgorithmHomeProject/internal/storage"
)

var ErrFromRemoteServer = errors.New("error from remote server")

type App struct {
	logger  Logger
	storage Storage
}

type DocumentSrc struct {
	Index int
	Url   string
	Title string
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
		Index: doc.Index, Url: doc.Url, Title: doc.Title,
	}
	return docApp
}

func documentSrcToStorage(docApp *DocumentSrc) storage.DocumentSource {
	docStor := storage.DocumentSource{
		Index: docApp.Index, Url: docApp.Url, Title: docApp.Title,
	}
	return docStor
}

func (a *App) getDocumentFromRemoteServer(url string) (string, string, error) {
	client := http.Client{}
	req, err := http.NewRequestWithContext(context.Background(), "GET", url, nil)
	if err != nil {
		return "", "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", "", ErrFromRemoteServer
	}
	doc, _ := goquery.NewDocumentFromReader(res.Body)
	doc.Find("script").Each(func(i int, el *goquery.Selection) {
		el.Remove()
	})
	return doc.Find("title").Text(), doc.Text(), nil
}

func (a *App) AddNewDocument(url string) ([]DocumentSrc, error) {
	docTitle, docText, err := a.getDocumentFromRemoteServer(url)
	if err != nil {
		return nil, err
	}
	_ = docText
	fmt.Println(docText)
	documents := make([]DocumentSrc, 0)
	err = a.storage.Add(storage.DocumentSource{Url: url, Title: docTitle, Data: docText})
	if err != nil {
		return nil, err
	}
	listDoc, err := a.storage.GetAllDocumentSource()
	if err != nil {
		return nil, err
	}
	for _, d := range listDoc {
		documents = append(documents, documentSrcFromStorage(&d))
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
