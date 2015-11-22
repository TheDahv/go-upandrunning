package nyt

import (
	"os"
	"testing"
)

func TestParseResponse(t *testing.T) {
	f, err := os.Open("./sample.json")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer f.Close()

	articles, err := parseResponse(f)

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	numArticles := len(articles)
	expectedArticles := 10

	if numArticles != expectedArticles {
		t.Errorf("Got %d articles, expected %d\n", numArticles, expectedArticles)
	}

	cases := []struct {
		Title string
		URL   string
	}{
		{
			"IN SUPPORT",
			"http://query.nytimes.com/gst/abstract.html?res=9C0DE4DF123BE633A25751C1A9659C946091D6CF",
		},
		{
			"Futility Lurked",
			"http://www.nytimes.com/1964/06/14/futility-lurked.html",
		},
		{
			"Remember Those Old-Time Beer Jingles?",
			"http://www.nytimes.com/1989/09/11/opinion/l-remember-those-old-time-beer-jingles-351089.html",
		},
		{
			"DID BEER SAVE THE BRITON?; To Us Intrenched in Prohibition, Science Says It Furnished Vitamins That Preserved Race The Beer of Old England. Canada's Self-Sacrifice. Cure for Scurvy. The Virtue of Beer. A Bare Living.",
			"http://query.nytimes.com/gst/abstract.html?res=9E04E4DC173AEE32A25752C2A9659C946195D6CF",
		},
		{
			"Is a Team of Clydesdales Next?; In Heady World of Microbreweries, Growth Is the Word",
			"http://www.nytimes.com/1994/10/25/nyregion/team-clydesdales-next-heady-world-microbreweries-growth-word.html",
		},
		{
			"Beer in Prague",
			"http://www.nytimes.com/1982/04/18/travel/l-beer-in-prague-067440.html",
		},
		{
			"MORE THAN ONE BREW.",
			"http://query.nytimes.com/gst/abstract.html?res=9407E7D81E3FE633A25750C1A9649C946394D6CF",
		},
		{
			"VARIETIES OF BEER.; Comparison of German and English Brews -- Effect of Tastes Upon Trade.",
			"http://query.nytimes.com/gst/abstract.html?res=9406EED91030E333A2575BC1A9649D94699ED7CF",
		},
		{
			"AN APPEAL FOR PURE BEER.",
			"http://query.nytimes.com/gst/abstract.html?res=980CE0D9133EE333A25750C2A9649D946197D6CF",
		},
		{
			"ADULTERATED LAGER BEER.; THE BRIDGE TOLLS.",
			"http://query.nytimes.com/gst/abstract.html?res=9B00E1DC1439E533A2575BC1A9609C94649FD7CF",
		},
	}

	for i, c := range cases {
		title := articles[i].Headline.Main
		url := articles[i].URL

		if title != c.Title {
			t.Errorf("Got title '%s', expected '%s'\n", title, c.Title)
		}

		if url != c.URL {
			t.Errorf("Got URL '%s', expected '%s'\n", url, c.URL)
		}
	}
}
