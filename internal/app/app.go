package app

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

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
	Title   string
	Context string
}

type WordInfo struct {
	IDDocument    int `json:"id_document"`
	PosInDocument int `json:"pos"`
}

type WordWithPos struct {
	Word string
	Pos  int
}

type Logger interface {
	Error(msg string)
}

type Storage interface {
	Add(document storage.Document) (int, error)
	Get(documentId int) (storage.Document, error)
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

func removeDuplicateStrings(s []WordWithPos) []WordWithPos {
	if len(s) < 1 {
		return s
	}

	sort.SliceStable(s, func(i, j int) bool {
		return s[i].Word < s[j].Word
	})

	prev := 1
	for curr := 1; curr < len(s); curr++ {
		if s[curr-1].Word != s[curr].Word {
			s[prev] = s[curr]
			prev++
		}
	}

	return s[:prev]
}

var asciiSpace = [256]uint8{'\t': 1, '\n': 1, '\v': 1, '\f': 1, '\r': 1, ' ': 1}

func fieldsFunc(s string, f func(rune) bool) []WordWithPos {
	// A span is used to record a slice of s of the form s[start:end].
	// The start index is inclusive and the end index is exclusive.
	type span struct {
		start int
		end   int
	}
	spans := make([]span, 0, 32)

	// Find the field start and end indices.
	// Doing this in a separate pass (rather than slicing the string s
	// and collecting the result substrings right away) is significantly
	// more efficient, possibly due to cache effects.
	start := -1 // valid span start if >= 0
	for end, rune := range s {
		if f(rune) {
			if start >= 0 {
				spans = append(spans, span{start, end})
				// Set start to a negative value.
				// Note: using -1 here consistently and reproducibly
				// slows down this code by a several percent on amd64.
				start = ^start
			}
		} else {
			if start < 0 {
				start = end
			}
		}
	}

	// Last field might end at EOF.
	if start >= 0 {
		spans = append(spans, span{start, len(s)})
	}

	// Create strings from recorded field indices.
	a := make([]WordWithPos, len(spans))
	for i, span := range spans {
		a[i] = WordWithPos{strings.ToLower(s[span.start:span.end]), span.start}
	}

	return a
}

func parseWordsFromText(s string) []WordWithPos {
	// First count the fields.
	// This is an exact count if s is ASCII, otherwise it is an approximation.
	n := 0
	wasSpace := 1
	// setBits is used to track which bits are set in the bytes of s.
	setBits := uint8(0)
	for i := 0; i < len(s); i++ {
		r := s[i]
		setBits |= r
		isSpace := int(asciiSpace[r])
		n += wasSpace & ^isSpace
		wasSpace = isSpace
	}

	if setBits >= utf8.RuneSelf {
		// Some runes in the input string are not ASCII.
		return fieldsFunc(s, unicode.IsSpace)
	}
	// ASCII fast path
	a := make([]WordWithPos, n)
	na := 0
	fieldStart := 0
	i := 0
	// Skip spaces in the front of the input.
	for i < len(s) && asciiSpace[s[i]] != 0 {
		i++
	}
	fieldStart = i
	for i < len(s) {
		if asciiSpace[s[i]] == 0 {
			i++
			continue
		}
		a[na] = WordWithPos{strings.ToLower(s[fieldStart:i]), fieldStart}
		na++
		i++
		// Skip spaces in between fields.
		for i < len(s) && asciiSpace[s[i]] != 0 {
			i++
		}
		fieldStart = i
	}
	if fieldStart < len(s) { // Last field might end at EOF.
		a[na] = WordWithPos{strings.ToLower(s[fieldStart:]), fieldStart}
	}
	return a
}

func textToSliceWord(text string) []WordWithPos {
	words := parseWordsFromText(text)
	words = removeDuplicateStrings(words)
	return words
}

func createAndSaveInvertIndex(storage *Storage, id int, text string) {
	words := textToSliceWord(text)
	for _, w := range words {
		wordInfoByte, err := (*storage).GetWordInfo(w.Word)
		if err != nil {
			panic(err) // TODO: add channel for return error
		}
		wil := make([]WordInfo, 0)
		if wordInfoByte != nil {
			json.Unmarshal(wordInfoByte, &wil)
		}
		wil = append(wil, WordInfo{id, w.Pos})
		wordInfoNewByte, err := json.Marshal(wil)
		if err != nil {
			panic(err) // TODO: add channel for return error
		}
		err = (*storage).UpdateWordInfo(w.Word, wordInfoNewByte)
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
	removePunctuation := func(r rune) rune {
		if strings.ContainsRune(".,:;!?[]()<>", r) {
			return ' '
		} else {
			return r
		}
	}
	docText = strings.Map(removePunctuation, docText)

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

const cFD = 100

type PositionInDoc struct {
	begin int
	end   int
}

func (a *App) Search(str string) ([]SearchResult, error) {
	result := make([]SearchResult, 0)
	words := textToSliceWord(str)
	mapIdRelevantDocuments := make(map[int]int)
	mapIdRelevantPosInDocument := make(map[int][]PositionInDoc)
	for _, word := range words {
		wordInfoByte, err := a.storage.GetWordInfo(word.Word)
		if err != nil {
			return nil, err
		}
		wil := make([]WordInfo, 0)
		if wordInfoByte != nil {
			json.Unmarshal(wordInfoByte, &wil)
		}
		for _, wordInfo := range wil {
			mapIdRelevantDocuments[wordInfo.IDDocument]++
			mapIdRelevantPosInDocument[wordInfo.IDDocument] =
				append(mapIdRelevantPosInDocument[wordInfo.IDDocument],
					PositionInDoc{wordInfo.PosInDocument, wordInfo.PosInDocument + len(word.Word)})
		}
	}
	indexMatch := 0
	for idDoc, countMatch := range mapIdRelevantDocuments {
		if countMatch != len(words) {
			continue
		}
		doc, err := a.storage.Get(idDoc)
		if err != nil {
			return nil, err
		}
		indexMatch++
		context := ""
		for _, pos := range mapIdRelevantPosInDocument[idDoc] {
			if pos.begin-cFD >= 0 && pos.end+cFD < len(doc.Data) {
				context += "..." + doc.Data[pos.begin-cFD:pos.end+cFD] + "..."
			} else if pos.begin-cFD >= 0 {
				context += "..." + doc.Data[pos.begin-cFD:] + "..."
			} else if pos.end+cFD < len(doc.Data) {
				context += "..." + doc.Data[:pos.end+cFD] + "..."
			}
		}
		result = append(result, SearchResult{indexMatch, doc.Url, doc.Title, context})

	}
	return result, nil
}

func New(logger Logger, storage Storage) *App {
	return &App{logger: logger, storage: storage}
}
