package libgen

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type LibgenBook struct {
	id        int
	authors   []string
	title     string
	publisher string
	year      int
	pages     int
	language  string
	size      string
	extension string
	mirrors   []string
}

func (b LibgenBook) Title() string    { return b.title }
func (b LibgenBook) Author() string   { return strings.Join(b.authors, ", ") }
func (b LibgenBook) Format() string   { return b.extension }
func (b LibgenBook) Size() string     { return b.size }
func (b LibgenBook) Language() string { return b.language }

func (b LibgenBook) findSupportedMirror() (*url.URL, error) {
	for _, mirror := range b.mirrors {
		u, err := url.Parse(mirror)
		if err != nil {
			continue
		}
		if u.Host == "library.lol" {
			return u, nil
		}
	}
	return nil, fmt.Errorf("no supported mirrors found")
}

func (b LibgenBook) DownloadLink() (string, error) {
	return getDownloadLink(b.mirrors)
}

type LibgenSearchProvider struct{}

func (l LibgenSearchProvider) Find(query string) ([]Book, error) {
	u, _ := url.Parse("http://libgen.rs/search.php")
	q := u.Query()
	q.Add("req", query)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err

	}

	body := findBody(doc)
	if body == nil {
		return nil, fmt.Errorf("parse error: cannot find body")
	}

	resultTable := eachChild(body, func(el *html.Node, _ int) *html.Node {
		if el.Type == html.ElementNode && el.Data == "table" {
			class, _ := getAttr(el, "class")
			if class == "c" {
				return el.FirstChild
			}
		}
		return nil
	})

	if resultTable == nil {
		return nil, fmt.Errorf("Parse error: cannot find results table")
	}

	var books []Book
	eachSiblingOfType(resultTable.FirstChild.NextSibling, "tr", func(row *html.Node, _ int) *html.Node {
		book := LibgenBook{}
		eachChildOfType(row, "td", func(td *html.Node, idx int) *html.Node {
			content := strings.TrimSpace(getText(td))
			switch idx {
			case 0: // ID
				book.id, _ = strconv.Atoi(getText(td))
			case 1: // Authors
				eachChildOfType(td, "a", func(authorLink *html.Node, _ int) *html.Node {
					name := strings.TrimSpace(getText(authorLink))
					if name == "" {
						return nil
					}
					book.authors = append(book.authors, name)
					return nil
				})
			case 2: // Title:
				book.title = content
			case 3:
				book.publisher = content
			case 4:
				book.year, _ = strconv.Atoi(content)
			case 5:
				book.pages, _ = strconv.Atoi(content)
			case 6:
				book.language = content
			case 7:
				book.size = content
			case 8:
				book.extension = content
			default:
				if idx >= 9 && content != "[edit]" && td.FirstChild != nil {
					href, ok := getAttr(td.FirstChild, "href")
					if ok {
						book.mirrors = append(book.mirrors, href)
					}
				}
			}
			return nil
		})
		books = append(books, book)
		return nil
	})

	return books, nil

}
