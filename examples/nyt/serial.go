package nyt

// This expands on the original example to demonstrate combining multiple
// searches into one list of responses. We're going to start with a serial
// implementation. Network calls won't start until the previous one is
// finished processing. We'll later improve upon this with a parallel solution.
//
// http://developer.nytimes.com/docs/read/article_search_api_v2

import (
	"bytes"
	"fmt"
	"io"
)

// RunSerial fetches articles for multiple search terms, one after the other.
func RunSerial(searchTerms []string) (io.Reader, error) {
	// Declaring variables this ways creates the 0-value of a type. The 0-value OMIT
	// of a slice is a nil pointer though, so we either need to initialize it OMIT
	// with `make`, or we can rely on `append` below OMIT
	var articles []Article

	for _, t := range searchTerms {
		a, err := findArticles(t)
		// Let's skip over individual errors for now and focus on the program OMIT
		if err != nil {
			fmt.Printf("Error searching for '%s': %s\n", t, err.Error())
			continue
		}
		// `append` adds the arguments to the destination slice. It grows the OMIT
		// underlying array and copies the elements when necessary. OMIT
		// Note the syntax for applying multiple arguments to variadic functions OMIT
		articles = append(articles, a...)
	}

	var rdr io.Reader
	if len(articles) == 0 {
		fmt.Println("Found no articles")
		return rdr, nil
	}

	var buf bytes.Buffer
	for _, article := range articles {
		fmt.Fprintf(&buf, "%s\n", article)
	}
	rdr = bytes.NewReader(buf.Bytes())
	return rdr, nil
}
