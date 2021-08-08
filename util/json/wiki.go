// Package json is a package for json files
package json

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	wikiTitle string
	wikiMsg   string
)

type wiki struct {
	Title string `json:"title"`
	Line  string `json:"line"`
}

func readWikiJSON() {
	wikiFile, err := os.Open("./asset/wiki.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer wikiFile.Close()
	var wikiData wiki
	byteValue, err := ioutil.ReadAll(wikiFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(byteValue, &wikiData)
	wikiTitle = "**" + wikiData.Title + "**"
	wikiMsg = wikiData.Line
}
