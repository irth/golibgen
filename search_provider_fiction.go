package libgen

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

type FictionBook struct {
	authors   []string
	series    string
	title     string
	language  string
	extension string
	size      string
	mirrors   []string
}

func (b FictionBook) Title() string  { return b.title }
func (b FictionBook) Author() string { return strings.Join(b.authors, ", ") }
func (b FictionBook) Format() string { return b.extension }

func (b FictionBook) DownloadLink() (string, error) {
	return getDownloadLink(b.mirrors)
}

type FictionSearchProvider struct{}

func (l FictionSearchProvider) Find(query string) ([]Book, error) {
	u, _ := url.Parse("http://libgen.rs/fiction/")
	q := u.Query()
	q.Add("q", query)
	u.RawQuery = q.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err

	}

	resultTable := findElement(doc, func(el *html.Node) bool {
		if el.Type == html.ElementNode && el.Data == "tbody" {
			class, _ := getAttr(el.Parent, "class")
			if class == "catalog" {
				return true
			}
		}
		return false
	})

	if resultTable == nil {
		return nil, fmt.Errorf("Parse error: cannot find results table")
	}

	var books []Book
	eachSiblingOfType(resultTable.FirstChild.NextSibling, "tr", func(row *html.Node, _ int) *html.Node {
		book := FictionBook{}
		eachChildOfType(row, "td", func(td *html.Node, idx int) *html.Node {
			content := strings.TrimSpace(getText(td))
			switch idx {
			case 0:
				ul := findElement(td, func(el *html.Node) bool { return el.Type == html.ElementNode && el.Data == "ul" })
				eachChildOfType(ul, "li", func(li *html.Node, _ int) *html.Node {
					book.authors = append(book.authors, getText(li))
					return nil
				})
			case 1:
				book.series = content
			case 2: // Title:
				book.title = content
				fmt.Println(content)
			case 3:
				book.language = content
			case 4:
				a := strings.Split(content, "/")
				if len(a) != 2 {
					book.extension = "unknown"
					book.size = "unknown"
				} else {
					book.extension = strings.ToLower(strings.TrimSpace(a[0]))
					book.size = strings.TrimSpace(a[1])
				}
			case 5:
				ul := findElement(td, func(el *html.Node) bool { return el.Type == html.ElementNode && el.Data == "ul" })
				eachChildOfType(ul, "li", func(li *html.Node, _ int) *html.Node {
					href, ok := getAttr(li.FirstChild, "href")
					if ok {
						book.mirrors = append(book.mirrors, href)
					}
					return nil
				})
			}
			return nil
		})
		books = append(books, book)
		return nil
	})

	return books, nil

}
