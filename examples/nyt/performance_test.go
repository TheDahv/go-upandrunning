package nyt

import (
	"io"
	"strconv"
	"testing"
	"time"
)

// Go's Benchmarking features in the testing package allow us to pinpoint slow
// parts of our program when subjected to a hot loop, as well as compare
// implementations for a given solution.
//
// Here, we introduce a ArticlesFetcher implementation that simulates network
// latency to demonstrate how our RunParallel solution reduces out IO costs.
//
// Run with: `go test -bench . ./...` to see results. On my Macbook Pro on Go
// 1.4, I get these results:
//
// 	$ go test -bench . ./...
// 	?       github.com/thedahv/go-upandrunning      [no test files]
// 	?       github.com/thedahv/go-upandrunning/examples/hello       [no test files]
// 	Reading key from ENV
//
// 	PASS
// 	BenchmarkRunSerial             2         525209466 ns/op
// 	BenchmarkRunParallel          20          57999243 ns/op
// 	ok      github.com/thedahv/go-upandrunning/examples/nyt 2.809s
//
// We can also control the concurrency with GOMAXPROCS to see how our program is
// affected by different concurrency levels (NOTE that Go 1.5 sets GOMAXPROCS to
// the number of CPU's on the machine, but earlier versions fo Go do not):
//
// 	$ GOMAXPROCS=8 go test -bench . ./...
// 	?       github.com/thedahv/go-upandrunning      [no test files]
// 	?       github.com/thedahv/go-upandrunning/examples/hello       [no test files]
// 	Reading key from ENV
//
// 	PASS
// 	BenchmarkRunSerial-8           2         531811734 ns/op
// 	BenchmarkRunParallel-8        30          55197976 ns/op
// 	ok      github.com/thedahv/go-upandrunning/examples/nyt 3.313s
//
// For more reading on performance and profiling, read this blog post:
// http://blog.golang.org/profiling-go-programs
//
// See if you can leverage Benchmark runs with Profiles to find areas to
// optimize the program.

// SlowFetcher simulates network latency while affording the predictability of a
// File Fetcher.
type SlowFetcher struct {
	fetcher FileFetcher
}

// Fetch introduces a small amount of wait time before delegating the call to
// the wrapped FileFetcher, returning the articles JSON bytes data inside.
func (sf SlowFetcher) Fetch(searchTerm string) (io.Reader, error) {
	time.Sleep(50 * time.Millisecond)
	return sf.fetcher.Fetch(searchTerm)
}

// NewSlowFetcher wraps a FileFetcher for the given path and returns the
// SlowFetcher.
func NewSlowFetcher(path string) (SlowFetcher, error) {
	var sf SlowFetcher

	ff, err := NewFileFetcher(path)
	if err != nil {
		return sf, err
	}

	return SlowFetcher{ff}, nil
}

var terms = func() []string {
	n := 10
	t := make([]string, n)

	for i := 0; i < n; i++ {
		t[i] = strconv.Itoa(i)
	}

	return t
}()

func BenchmarkRunSerial(t *testing.B) {
	sf, err := NewSlowFetcher("./sample.json")

	if err != nil {
		t.Error()
		t.FailNow()
	}

	for i := 0; i < t.N; i++ {
		RunSerial(sf, terms)
	}
}

func BenchmarkRunParallel(t *testing.B) {
	sf, err := NewSlowFetcher("./sample.json")

	if err != nil {
		t.Error()
		t.FailNow()
	}

	for i := 0; i < t.N; i++ {
		RunParallel(sf, terms)
	}
}
