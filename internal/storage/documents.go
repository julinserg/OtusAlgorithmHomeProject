package storage

type DocumentSource struct {
	Index int    `db:"id"`
	Url   string `db:"url"`
	Title string `db:"title"`
	Data  string `db:"data"`
}
