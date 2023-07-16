package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
	"github.com/julinserg/OtusAlgorithmHomeProject/internal/storage"
)

var ErrFromRemoteServer = errors.New("error from remote server")

type App struct {
	logger  Logger
	storage Storage
}

type DocumentSrc struct {
	ID        int
	SeqNumber int
	Url       string
	Title     string
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

func documentSrcFromStorage(doc *storage.DocumentSource, index int) DocumentSrc {
	docApp := DocumentSrc{
		ID: doc.ID, SeqNumber: index + 1, Url: doc.Url, Title: doc.Title,
	}
	return docApp
}

func documentSrcToStorage(docApp *DocumentSrc) storage.DocumentSource {
	docStor := storage.DocumentSource{
		ID: docApp.ID, Url: docApp.Url, Title: docApp.Title,
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
	regex, err := regexp.Compile("\\s+\n")
	if err != nil {
		return "", "", err
	}
	title := ""
	doc.Find("title").EachWithBreak(func(i int, el *goquery.Selection) bool {
		title = el.Text()
		return false
	})
	text := regex.ReplaceAllString(doc.Text(), "\n")
	return title, text, nil
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
	for index, d := range listDoc {
		documents = append(documents, documentSrcFromStorage(&d, index))
	}
	return documents, nil
}

func (a *App) GetAllDocument() ([]DocumentSrc, error) {
	documents := make([]DocumentSrc, 0)
	listDoc, err := a.storage.GetAllDocumentSource()
	if err != nil {
		return nil, err
	}
	for index, d := range listDoc {
		documents = append(documents, documentSrcFromStorage(&d, index))
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
