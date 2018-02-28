package wotd

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// LongmanURL is an URL for getting word of the day.
const LongmanURL = "https://www.ldoceonline.com/"

func fetchPage(url string) (*goquery.Document, error) {
	doc, err := goquery.NewDocument(url)
	return doc, err
}

func fetchWOTD(mainURL string) (*WOTD, error) {
	mainPage, err := fetchPage(mainURL)
	if err != nil {
		return nil, err
	}
	wordURL, err := parseWordURL(mainPage)
	if err != nil {
		return nil, err
	}
	wordPage, err := fetchPage(wordURL)
	if err != nil {
		return nil, err
	}
	wotd, err := parseWOTD(wordPage)
	return wotd, err
}

func getCachePath() string {
	date := time.Now().Local().Format("2006-01-02")
	return path.Join(os.TempDir(), fmt.Sprintf("wotd-%s.json", date))
}

func getCachedWord(filepath string) (*WOTD, error) {
	wotd := &WOTD{}
	fi, err := os.Open(filepath)
	defer fi.Close()
	if err != nil {
		return wotd, err
	}
	err = wotd.Load(fi)
	return wotd, err
}

func saveCachedWord(wotd *WOTD, filepath string) error {
	fo, err := os.Create(filepath)
	defer fo.Close()
	err = wotd.Dump(fo)
	return err
}

// GetWOTD gets the word of the day and its definition
// from the Longman dictionary and caches it.
func GetWOTD() (*WOTD, error) {
	var wotd *WOTD
	cachePath := getCachePath()
	wotd, err := getCachedWord(cachePath)
	if err != nil {
		wotd, err = fetchWOTD(LongmanURL)
		if err != nil {
			return wotd, err
		}
		err = saveCachedWord(wotd, cachePath)
	}
	return wotd, err
}
