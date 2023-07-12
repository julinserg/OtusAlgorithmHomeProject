package storage

type DocumentSource struct {
	Index int    `db:"id"`
	Url   string `db:"url"`
}
