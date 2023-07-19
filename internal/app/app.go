package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/julinserg/OtusAlgorithmHomeProject/internal/storage"
)

var ErrFromRemoteServer = errors.New("error from remote server")

type App struct {
	logger  Logger
	storage Storage
}

type Document struct {
	ID        int
	SeqNumber int
	Url       string
	Title     string
}

type SearchResult struct {
	Index   int
	Url     string
	Context string
}

type WordInfo struct {
	IDDocument    int `json:"id_document"`
	PosInDocument int `json:"pos"`
}

type Logger interface {
	Error(msg string)
}

type Storage interface {
	Add(document storage.Document) (int, error)
	GetAllDocumentSource() ([]storage.Document, error)
	GetWordInfo(word string) ([]byte, error)
	UpdateWordInfo(word string, wordInfo []byte) error
}

func documentSrcFromStorage(doc *storage.Document, index int) Document {
	docApp := Document{
		ID: doc.ID, SeqNumber: index + 1, Url: doc.Url, Title: doc.Title,
	}
	return docApp
}

func documentSrcToStorage(docApp *Document) storage.Document {
	docStor := storage.Document{
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

func removeDuplicateStrings(s []string) []string {
	if len(s) < 1 {
		return s
	}

	sort.Strings(s)
	prev := 1
	for curr := 1; curr < len(s); curr++ {
		if s[curr-1] != s[curr] {
			s[prev] = s[curr]
			prev++
		}
	}

	return s[:prev]
}

func toLowerStrings(s []string) []string {
	if len(s) < 1 {
		return s
	}
	for curr := 1; curr < len(s); curr++ {
		s[curr] = strings.ToLower(s[curr])
	}
	return s
}

func createAndSaveInvertIndex(storage *Storage, id int, text string) {
	removePunctuation := func(r rune) rune {
		if strings.ContainsRune(".,:;!?[]()<>", r) {
			return -1
		} else {
			return r
		}
	}

	s := strings.Map(removePunctuation, text)
	words := strings.Fields(s)
	words = toLowerStrings(words)
	words = removeDuplicateStrings(words)
	for _, w := range words {
		fmt.Println(w)
		wordInfoByte, err := (*storage).GetWordInfo(w)
		if err != nil {
			panic(err) // TODO: add channel for return error
		}
		wil := make([]WordInfo, 0)
		if wordInfoByte != nil {
			json.Unmarshal(wordInfoByte, &wil)
		}
		wil = append(wil, WordInfo{id, 0})
		wordInfoNewByte, err := json.Marshal(wil)
		if err != nil {
			panic(err) // TODO: add channel for return error
		}
		err = (*storage).UpdateWordInfo(w, wordInfoNewByte)
		if err != nil {
			panic(err) // TODO: add channel for return error
		}
	}
}

func (a *App) AddNewDocument(url string) ([]Document, error) {
	docTitle, docText, err := a.getDocumentFromRemoteServer(url)
	if err != nil {
		return nil, err
	}

	documents := make([]Document, 0)
	id, err := a.storage.Add(storage.Document{Url: url, Title: docTitle, Data: docText})
	if err != nil {
		return nil, err
	}

	go createAndSaveInvertIndex(&a.storage, id, docText)

	listDoc, err := a.storage.GetAllDocumentSource()
	if err != nil {
		return nil, err
	}
	for index, d := range listDoc {
		documents = append(documents, documentSrcFromStorage(&d, index))
	}
	return documents, nil
}

func (a *App) GetAllDocument() ([]Document, error) {
	documents := make([]Document, 0)
	listDoc, err := a.storage.GetAllDocumentSource()
	if err != nil {
		return nil, err
	}
	for index, d := range listDoc {
		documents = append(documents, documentSrcFromStorage(&d, index))
	}
	return documents, nil
}

func (a *App) Search(str string) ([]SearchResult, error) {
	documents := make([]SearchResult, 0)
	documents = append(documents, SearchResult{
		1,
		"https://stackoverflow.com/questions/9523927/how-to-stop-table-from-resizing-when-contents-grow",
		"I have a table, the cells of which are filled with picture"})
	documents = append(documents, SearchResult{
		2,
		"https://stackoverflow.com/questions/21019302/html-button-layout-positioning",
		"Even i didn't get what exactly you want. but for an image sourrounded by buttons try this code"})
	return documents, nil
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}
