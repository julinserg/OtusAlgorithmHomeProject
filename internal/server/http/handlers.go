package internalhttp

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/julinserg/OtusAlgorithmHomeProject/internal/app"
)

var landingFormTmpl = `
<html>
	<head>
		<style>
			h1 {
				text-align: center;
			}	
		</style>
	</head>
    <body>
	<h1>Mini Search</h1>
	<table align="center">
	<tr>
		<td style="text-align: center"><form action="/add" method="post">
			Add: <input type="text" name="add">
			<input type="submit" value="Add">
		</form></td>
		<td style="text-align: center"><form action="/search" method="post">
			Search: <input type="text" name="search">
			<input type="submit" value="Search">
		</form></td>
	</tr>	     
	</table>    
	</body>
</html>
`

var resultFormTmpl = `
<html>
	<head>
		<style>
			h1 {
				text-align: center;
			}	
		</style>
	</head>
    <body>
	<h1>Mini Search</h1>
	<table align="center">
	<tr>
		<td style="text-align: center"><form action="/add" method="post">
			Add: <input type="text" name="add">
			<input type="submit" value="Add">
		</form></td>
		<td style="text-align: center"><form action="/search" method="post">
			Search: <input type="text" name="search">
			<input type="submit" value="Search">
		</form></td>
	</tr>
	<tr>	
	<td style="text-align: center">	
	<table>
	{{ range .ItemsSource}}
		<tr>
			<td>{{ .Index }}</td>	
			<td>{{ .Url }}</td>			
		</tr>		
	{{ end}}
	</table>
	</td>
	<td style="text-align: center">	
	<table>
	{{ range .ItemsResult}}
		<tr>
			<td>{{ .Index }}</td>
			<td>{{ .Url }}</td>		
		</tr>	
		<tr>
			<td></td>
			<td>{{ .Context }}</td>
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
	ItemsSource []app.DocumentSrc
	ItemsResult []app.DocumentSearch
}

type minisearchHandler struct {
	logger Logger
	app    Application
	data   Data
}

func (ph *minisearchHandler) landingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(landingFormTmpl))
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
	t := template.Must(template.New("result").Parse(resultFormTmpl))
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
	t := template.Must(template.New("result").Parse(resultFormTmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, ph.data); err != nil {
		panic(err)
	}
	s := buf.String()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s))
}
