package wotd

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/fatih/color"
)

// WordSense describes a word meaning with examples.
type WordSense struct {
	Definition string
	Examples   []string
}

// DictEntry describes an entry with several meanings.
type DictEntry struct {
	PartOfSpeech string
	Head         string
	Senses       []WordSense
}

// WOTD describes word and its definitions with examples.
type WOTD struct {
	Word    string
	Entries []DictEntry
}

// String generates string representation of WOTD.
func (wotd *WOTD) String() string {
	result := fmt.Sprintf("Word Of The Day: %s\n", color.YellowString(wotd.Word))
	for _, entry := range wotd.Entries {
		result += fmt.Sprintf("%s (%s)\n", color.YellowString(entry.Head), color.GreenString(entry.PartOfSpeech))
		for i, sense := range entry.Senses {
			result += fmt.Sprintf("\t%v. %s\n", i+1, sense.Definition)
			for _, example := range sense.Examples {
				result += fmt.Sprintf("\t\t* %s\n", color.CyanString(example))
			}
		}
	}
	return result
}

// Dump serializes WOTD into JSON and writes to io.Writer.
func (wotd *WOTD) Dump(w io.Writer) error {
	data, err := json.Marshal(wotd)
	if err != nil {
		return err
	}
	_, err = w.Write(data)
	return err
}

// Load deserializes data from io.Reader.
func (wotd *WOTD) Load(r io.Reader) error {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read cache file: %s", err)
	}
	err = json.Unmarshal(data, wotd)
	return err
}
