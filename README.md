# golibgen

## Installation
```
go get github.com/irth/golibgen
```

## Usage example

```golang
package main

import (
	"fmt"

	libgen "github.com/irth/golibgen"
)

func main() {
	sp := libgen.FictionSearchProvider{}
	books, _ := sp.Find("query")
	fmt.Println(books[0].DownloadLink())
}
```
