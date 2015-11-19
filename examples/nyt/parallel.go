package nyt

import (
	"bytes"
	"fmt"
	"io"
	"sync"
)

// Most programmers--especially Node programmers--know blocking, synchronous
// serial code is not the best approach to IO-intensive programs.
//
// Go lets us take advantage of lightweight, thread-like goroutines to run
// procedures in parallel. The other big win is Go lets us take advantage of all
// our CPU's to do work concurrently.
//
// NOTE: concurrency primitives aren't free, and you'll need to do your
// profiling to determine if the work across your CPU's outweighs the cost of
// scheduling and coordinating.
//
// There is a good talk about the difference here:
// http://blog.golang.org/concurrency-is-not-parallelism
//
// Other good reading:
// - https://blog.golang.org/pipelines
// - https://talks.golang.org/2012/concurrency.slide

// RunParallel fetches articles for multiple search terms in parallel and merges
// the results into a single payload.
func RunParallel(searchTerms []string) (io.Reader, error) {
	// setup start OMIT
	var (
		articles []Article
		rdr      io.Reader
		buf      bytes.Buffer
	)

	// A Waitgroup allows our goroutines to signal when they are completed OMIT
	// working, and for the main thread to wait until all are done before OMIT
	// proceeding. OMIT
	var wg sync.WaitGroup
	// We already know how many goroutines we will run, so add that now OMIT
	wg.Add(len(searchTerms))

	// results lets each goroutine send Articles back to the main thread OMIT
	// Rather than sharing access to `articles`, we share by communicating OMIT
	results := make(chan []Article)
	// errors allows each goroutine to signal an error if there is one OMIT
	errors := make(chan error)
	// We can use this buffered channel to signal all work completed OMIT
	done := make(chan bool, 1)

	// Note, channels aren't GC'd and you need to clean them yourself OMIT
	defer close(results)
	defer close(errors)
	defer close(done)

	// Fire off a search in a separate goroutine for each search term OMIT
	// NOTE: If we wanted to limit how fast we request our upstream source, we OMIT
	// could use a "worker pattern": OMIT
	for _, t := range searchTerms {
		go queryHelper(t, &wg, results, errors)
	}

	// We can run anonymous functions in their own goroutines as well OMIT
	go func() {
		// This will block until all goroutines signal they are done OMIT
		wg.Wait()
		// Now we can tell the "select" loop below that we are finished OMIT
		done <- true
	}()
	// setup end OMIT

	// Go's single loop construct can also serve as an infinite loop
	// For better or worse, Go also has labels. Note how we use that to
	// break out of the infinite loop from the select code.

	// consume start OMIT
Loop:
	for {
		select {
		case arts := <-results:
			articles = append(articles, arts...)
		case err := <-errors:
			fmt.Printf("Error in search: %s\n", err.Error())
		case <-done:
			break Loop
		}
	}

	// All the goroutines are done and we're back in single-threaded land OMIT
	if len(articles) == 0 {
		fmt.Println("Found no articles")
		return rdr, nil
	}

	for _, article := range articles {
		fmt.Fprintf(&buf, "%s\n", article)
	}

	rdr = bytes.NewReader(buf.Bytes())

	return rdr, nil
	// consume end OMIT
}

// queryHelper wraps our findArticles call with the WaitGroup synchronization
func queryHelper(term string, wg *sync.WaitGroup,
	results chan []Article, errors chan error) {

	articles, err := findArticles(term)
	if err != nil {
		errors <- err
	} else {
		results <- articles
	}

	wg.Done()
}
