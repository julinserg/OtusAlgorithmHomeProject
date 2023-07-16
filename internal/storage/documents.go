package storage

type DocumentSource struct {
	ID    int    `db:"id"`
	Url   string `db:"url"`
	Title string `db:"title"`
	Data  string `db:"data"`
}
