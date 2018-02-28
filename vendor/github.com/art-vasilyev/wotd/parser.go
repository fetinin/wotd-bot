package wotd

import (
	"errors"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func parseWordURL(doc *goquery.Document) (string, error) {
	var wordURL string
	wordURL, exist := doc.Find("#wotd .title_entry a").Attr("href")
	if exist != true {
		return wordURL, errors.New("unable to find WOTD link")
	}
	return wordURL, nil
}

func parseWOTD(doc *goquery.Document) (*WOTD, error) {
	wotd := WOTD{Word: strings.TrimSpace(doc.Find(".pagetitle").Text())}
	doc.Find(".dictionary .dictentry").Each(func(i int, s *goquery.Selection) {
		entry := DictEntry{
			Head:         strings.TrimSpace(s.Find(".Head span").First().Text()),
			PartOfSpeech: strings.TrimSpace(s.Find(".Head .POS").Text())}

		s.Find(".dictlink .Sense").Each(func(i int, s *goquery.Selection) {
			definition := strings.TrimSpace(s.Find(".DEF").Text())
			if definition == "" {
				return
			}
			sense := WordSense{Definition: definition}
			s.Find(".EXAMPLE").Each(func(i int, s *goquery.Selection) {
				sense.Examples = append(sense.Examples, strings.TrimSpace(s.Text()))
			})
			entry.Senses = append(entry.Senses, sense)
		})
		wotd.Entries = append(wotd.Entries, entry)
	})
	return &wotd, nil
}
