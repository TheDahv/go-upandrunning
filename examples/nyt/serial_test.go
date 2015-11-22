package nyt

import (
	"bufio"
	"testing"
)

func TestRunSerial(t *testing.T) {
	fetcher, err := NewFileFetcher("./sample.json")

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	rdr, err := RunSerial(fetcher, []string{"foo", "bar"})
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	scanner := bufio.NewScanner(rdr)
	actualArticles := 0
	expectedArticles := 20

	for scanner.Scan() {
		actualArticles++
	}

	if actualArticles != expectedArticles {
		t.Errorf("Got %d article results, expected %d\n",
			actualArticles, expectedArticles)
	}
}
