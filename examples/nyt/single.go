package nyt

// Here we have a simple program that takes a single article topic
// as input and builds a query to the New York Times Articles Search API.
//
// It then parses the results and writes the date, headline, and URL as a
// bytestream that can be consumed as an io.Reader.
//
// This is meant not to demonstrate high-performance approaches or low-level
// programming, but to demonstrate a common API consumer operation familiar
// to programmers from other languages.

import (
	"bytes"
	"fmt"
	"io"
)

// RunSingle fetches articles for a single search term from the NYT API.
func RunSingle(fetcher ArticlesFetcher, searchTerms []string) (io.Reader, error) {
	searchTerm := searchTerms[0]
	articles, err := findArticles(fetcher, searchTerm)

	if err != nil {
		fmt.Printf("Error finding articles: %s\n", err.Error())
		return nil, err
	}

	if len(articles) == 0 {
		fmt.Printf("Found no articles related to '%s'\n", searchTerm)
		return nil, err
	}

	// Collect results in a buffer before returning them wrapped in an io.Reader OMIT
	var buf bytes.Buffer
	// Introducing `range`. When ranging over arrays or slices, you get OMIT
	// an index and a value copy back. OMIT
	// NOTE: `article` is a copy and not a reference! OMIT
	for _, article := range articles {
		fmt.Fprintf(&buf, "%s\n", article)
	}

	// We have a buffer of data for our results, but we need to return an OMIT
	// interface. Specifically, the io.Reader interface describes anything that OMIT
	// can be read from with a Read() method. Keeping it less OMIT
	// implementation-specifics allows us to write decoupled, testable code. OMIT
	rdr := bytes.NewReader(buf.Bytes())
	return rdr, nil
}
