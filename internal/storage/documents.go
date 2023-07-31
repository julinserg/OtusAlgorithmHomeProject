package storage

type Document struct {
	ID    int    `db:"id"`
	URL   string `db:"url"`
	Title string `db:"title"`
	Data  string `db:"data"`
}
