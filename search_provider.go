package libgen

type SearchProvider interface {
	Find(query string) ([]Book, error)
}

type Book interface {
	Title() string
	Author() string
	Format() string
	Size() string
	Language() string
	DownloadLink() (string, error)
}
