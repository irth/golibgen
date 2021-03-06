package libgen

import (
	"strings"

	"golang.org/x/net/html"
)

func eachSibling(firstEl *html.Node, f func(*html.Node, int) *html.Node) *html.Node {
	idx := 0
	for el := firstEl; el != nil; el = el.NextSibling {
		if el.Type != html.ElementNode {
			continue
		}
		ret := f(el, idx)
		if ret != nil {
			return ret
		}
		idx++
	}
	return nil
}

func eachChild(el *html.Node, f func(*html.Node, int) *html.Node) *html.Node {
	if el == nil {
		return nil
	}
	return eachSibling(el.FirstChild, f)
}

func eachSiblingOfType(firstEl *html.Node, typ string, f func(*html.Node, int) *html.Node) *html.Node {
	return eachSibling(firstEl, func(el *html.Node, idx int) *html.Node {
		if el.Data != typ {
			return nil
		}
		return f(el, idx)
	})
}

func eachChildOfType(el *html.Node, typ string, f func(*html.Node, int) *html.Node) *html.Node {
	if el == nil {
		return nil
	}
	return eachSiblingOfType(el.FirstChild, typ, f)
}

func getText(el *html.Node) string {
	if el == nil {
		return ""
	}

	if el.Type == html.TextNode {
		return el.Data
	}

	if el.Type == html.ElementNode && el.Data == "br" {
		return ", "
	}

	children := []string{}
	for ch := el.FirstChild; ch != nil; ch = ch.NextSibling {
		children = append(children, getText(ch))
	}

	return strings.Join(children, "")

}

func findBody(n *html.Node) *html.Node {
	return findElement(n, func(el *html.Node) bool { return el.Type == html.ElementNode && el.Data == "body" })
}

func findElement(n *html.Node, predicate func(el *html.Node) bool) *html.Node {
	if predicate(n) {
		return n
	}
	var result *html.Node = nil
	return eachChild(n, func(c *html.Node, _ int) *html.Node {
		result = findElement(c, predicate)
		if result != nil {
			return result
		}
		return nil
	})
	return nil
}

func getAttr(n *html.Node, key string) (string, bool) {
	if n == nil {
		return "", false
	}
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val, true
		}
	}
	return "", false
}
