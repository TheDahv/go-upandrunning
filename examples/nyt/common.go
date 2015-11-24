package nyt

import (
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
	return fmt.Sprintf("(%v) '%s' - %s\n",
		a.Date.Format("02 Jan 2006"), a.Headline.Main, a.URL)
}

// Runnable is an example function that has a common signature
type Runnable func([]string) (io.Reader, error)

func findArticles(searchTerm string) ([]Article, error) {
	if key == "" {
		return nil, errors.New("API key not available")
	}

	resp, err := http.Get(fmt.Sprintf("%s.%s?q=%s&fq=%s&api-key=%s",
		urlBase, respFmt, searchTerm, docType, key))

	// Introducing "defer". This will run when the current scope returns 		OMIT
	// Defers can call `resp.Close()` directly, but we can also register 		OMIT
	// anonymous functions. We normally wouldn't do it this way, but we're 	OMIT
	// bailing early on error, so we must register the defer now. 					OMIT
	// 																																			OMIT
	// NOTE: Failing to close the response Body is a leak. Don't do it. 		OMIT
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

	// We're just going to choose to discard the response body ready for OMIT
	// brevity																													 OMIT
	// NOTE: This allocates the entire body into a buffer. You might not OMIT
	// need to do that if handling it as an io.Reader suffices.					 OMIT
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Send the raw bytes to `parseResponse` and return the []Articles OMIT
	return parseResponse(body)
}

// parseResponse marshals raw JSON bytes into a slice of Article structs.
// If there is a parsing error, that is returned with an empty articles list.
func parseResponse(data []byte) ([]Article, error) {
	var response ArticleResponse
	err := json.Unmarshal(data, &response)

	if err != nil {
		return []Article{}, err
	}
	return response.Response.Articles, nil
}
