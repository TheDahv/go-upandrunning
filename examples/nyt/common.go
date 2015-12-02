package nyt

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// Constants and variables can be declared outside of functions.
// Constants behave slightly differently than standard typed variables.
// You can read more here: https://blog.golang.org/constants
const (
	urlBase = "http://api.nytimes.com/svc/search/v2/articlesearch"
	respFmt = "json"
	docType = "article"
)

var key string

// Init is a special function to initialize a package. Any file in a package can
// have one, and a package can have multiple init() calls. This is useful for
// setting up values that will be used across the entire package.
func init() {
	// There are a few ways to handle secret values in our code:
	// from arguments, environment variables, or injected into our code
	// in the build process (see run.sh). Let's go with the last 2 options.
	if key == "" {
		fmt.Println("Reading key from ENV")
		key = os.Getenv("API_KEY")

		if key == "" {
			fmt.Println("No key found in environment variable API_KEY")
			return
		}
	} else {
		fmt.Println("Picked up key from build!")
	}
	fmt.Println()
}

// ArticleResponse models the API response itself. This demonstrates how we
// can use nested anonymous structs to model nested types. It's pretty ugly,
// so we can declare a specific type for each Article to keep things neat.
type ArticleResponse struct {
	Response struct {
		// Note that struct field names don't have to match JSON keys
		Articles []Article `json:"docs"`
	} `json:"response"`
}

// Article models a single article returned by the API
type Article struct {
	Headline struct {
		Main string `json:"main"`
	} `json:"headline"`
	URL  string    `json:"web_url"`
	Date time.Time `json:"pub_date"`
}

func (a Article) String() string {
	return fmt.Sprintf("(%v) '%s' - %s",
		a.Date.Format("02 Jan 2006"), a.Headline.Main, a.URL)
}

// Runnable is an example function that has a common signature
type Runnable func(ArticlesFetcher, []string) (io.Reader, error)

// ArticlesFetcher describes the behavior of something that can get
// articles from a resource.
type ArticlesFetcher interface {
	// Fetch returns the JSON payload for articles data wrapped in an io.Reader.
	Fetch(string) (io.Reader, error)
}

// NetworkFetcher implements the ArticlesFetcher interface by issuing HTTP
// requests against the NYT API for a given search term.
type NetworkFetcher struct {
	Key string
}

// NewNetworkFetcher returns a fetcher with the initialized API key. If none can
// be found from the environment, it returns an error.
func NewNetworkFetcher() (NetworkFetcher, error) {
	var nf NetworkFetcher
	if key == "" {
		return nf, errors.New("no API key found")
	}

	nf.Key = key
	return nf, nil
}

// Fetch creates the HTTP request for the search term and returns an io.Reader
// for the raw JSON response.
func (nf NetworkFetcher) Fetch(searchTerm string) (io.Reader, error) {
	resp, err := http.Get(fmt.Sprintf("%s.%s?q=%s&fq=%s&api-key=%s",
		urlBase, respFmt, searchTerm, docType, nf.Key))

	// Introducing "defer". This will run when the current scope returns 	 OMIT
	// Defers can call `resp.Close()` directly, but we can also register 	 OMIT
	// anonymous functions. We normally wouldn't do it this way, but we're OMIT
	// bailing early on error, so we must register the defer now. 				 OMIT
	// 																																		 OMIT
	// NOTE: Failing to close the response Body is a leak. Don't do it. 	 OMIT
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	// Infamous Go error-handling style. You're forced to handle OMIT
	// or ignore ALL errors.																		 OMIT
	if err != nil {
		fmt.Printf("Error issuing API request: %s\n", err.Error())
		return nil, err
	}

	// We're just going to choose to discard the response err ready for  OMIT
	// brevity. There are other elegant options like io.Pipe, as well    OMIT
	// NOTE: This allocates the entire body into a buffer. You might not OMIT
	// need to do that if handling it as an io.Reader suffices.					 OMIT
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	rdr := bytes.NewReader(body)
	return rdr, nil
}

// FileFetcher implements the ArticlesFetcher interface by reutrning data from a
// file on disk.
type FileFetcher struct {
	// data is an internally-managed byte array containing the data in the file
	data []byte
}

// NewFileFetcher attempts to wrap a FileFetcher around a file resource at path.
// It returns an error if there is a problem opening that file.
func NewFileFetcher(path string) (FileFetcher, error) {
	var ff FileFetcher

	f, err := os.Open(path)
	if err != nil {
		return ff, err
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return ff, err
	}

	ff.data = data

	return ff, nil
}

// Fetch loads articles from the JSON data that was extracted from
// a file on disk.
func (ff FileFetcher) Fetch(searchTerm string) (io.Reader, error) {
	if len(ff.data) == 0 {
		return nil, errors.New("FileFetcher file data not available")
	}

	d := make([]byte, len(ff.data))
	copy(d, ff.data)

	return bytes.NewReader(d), nil
}

func findArticles(fetcher ArticlesFetcher, searchTerm string) ([]Article, error) {
	data, err := fetcher.Fetch(searchTerm)
	if err != nil {
		return nil, err
	}

	// Send the raw bytes to `parseResponse` and return the []Articles OMIT
	return parseResponse(data)
}

// parseResponse marshals raw JSON bytes into a slice of Article structs.
// If there is a parsing error, that is returned with an empty articles list.
func parseResponse(data io.Reader) ([]Article, error) {
	var response ArticleResponse

	b, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &response)

	if err != nil {
		return []Article{}, err
	}
	return response.Response.Articles, nil
}
