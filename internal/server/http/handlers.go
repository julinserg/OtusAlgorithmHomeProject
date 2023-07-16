package internalhttp

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/julinserg/OtusAlgorithmHomeProject/internal/app"
)

var htmlFormTmpl = `
<html>
	<head>	
		<style>
		.ellipsis {
			max-width: 40px;
			text-overflow: ellipsis;
			overflow: hidden;
			white-space: nowrap;
		}
		</style>
	</head>	
    <body>
	<h1 style="text-align: center">Mini Search</h1>
	<table align="left" width="1200" cellpadding="10">
	<tr>
		<td style="text-align: left; vertical-align: top" width="50%"><form action="/add" method="post">
			Add: <input type="text" name="add" size="45">
			<input type="submit" value="Add">
		</form></td>
		<td style="text-align: left; vertical-align: top" width="50%"><form action="/search" method="post">
			Search: <input type="text" name="search" size="45">
			<input type="submit" value="Search">
		</form></td>	
	</tr>
	<tr>	
	<td style="text-align: left; vertical-align: top" width="50%">	
	<table>
	{{ range .ItemsSource}}
		<tr>
			<td width="10%">{{ .SeqNumber }}</td>	
			<td width="90%" class="ellipsis"><a href="{{ .Url }}">{{ .Url }}</a></td>			
		</tr>
		<tr>
			<td width="10%"></td>
			<td width="90%">{{ .Title }}</td>
		</tr>				
	{{ end}}
	</table>
	</td>
	<td style="text-align: left; vertical-align: top" width="50%">	
	<table>
	{{ range .ItemsResult}}
		<tr>
			<td width="10%">{{ .Index }}</td>
			<td width="90%" class="ellipsis"><a href="{{ .Url }}">{{ .Url }}</a></td>		
		</tr>	
		<tr>
			<td width="10%"></td>
			<td width="90%">{{ .Context }}</td>
		</tr>
	{{ end}}	
	</table>
	</td>	
	</tr>  	      
	</table>    
	</body>
</html>
`

type Data struct {
	ItemsSource []app.Document
	ItemsResult []app.SearchResult
}

type minisearchHandler struct {
	logger Logger
	app    Application
	data   Data
}

func (ph *minisearchHandler) landingHandler(w http.ResponseWriter, r *http.Request) {
	listDoc, err := ph.app.GetAllDocument()
	if err != nil {
		panic(err)
	}
	ph.data.ItemsSource = nil
	for _, doc := range listDoc {
		ph.data.ItemsSource = append(ph.data.ItemsSource, doc)
	}
	w.WriteHeader(http.StatusOK)
	t := template.Must(template.New("result").Parse(htmlFormTmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, ph.data); err != nil {
		panic(err)
	}
	s := buf.String()
	w.Write([]byte(s))
}

func (ph *minisearchHandler) searchHandler(w http.ResponseWriter, r *http.Request) {

	searchString := r.FormValue("search")
	listResultSearch, err := ph.app.Search(searchString)
	if err != nil {
		panic(err)
	}
	ph.data.ItemsResult = nil
	for _, doc := range listResultSearch {
		ph.data.ItemsResult = append(ph.data.ItemsResult, doc)
	}
	t := template.Must(template.New("result").Parse(htmlFormTmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, ph.data); err != nil {
		panic(err)
	}
	s := buf.String()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}

func (ph *minisearchHandler) addHandler(w http.ResponseWriter, r *http.Request) {

	urlDoc := r.FormValue("add")
	listDoc, err := ph.app.AddNewDocument(urlDoc)
	if err != nil {
		panic(err)
	}
	ph.data.ItemsSource = nil
	for _, doc := range listDoc {
		ph.data.ItemsSource = append(ph.data.ItemsSource, doc)
	}
	t := template.Must(template.New("result").Parse(htmlFormTmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, ph.data); err != nil {
		panic(err)
	}
	s := buf.String()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}
