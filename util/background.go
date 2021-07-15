// Package util is a package for json files and database.
package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	backgroundTitle string
	backgroundMsg   string
)

type background struct {
	Title string   `json:"title"`
	Line  []bgLine `json:"line"`
}

type bgLine struct {
	Pre  string `json:"pre"`
	Cmd  string `json:"cmd"`
	Post string `json:"post"`
}

// background.json 파일 읽어서 "ㅁ게임배경" 실행시 출력할 데이터 세팅
func readBackgroundJSON(prefix string) {
	jsonFile, err := os.Open("./asset/background.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer jsonFile.Close()
	var bg background
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(byteValue, &bg)

	backgroundTitle = "**" + bg.Title + "**"
	backgroundMsg = ""
	for i := 0; i < len(bg.Line); i++ {
		if len(bg.Line[i].Pre) > 0 {
			backgroundMsg += bg.Line[i].Pre
		}
		if len(bg.Line[i].Cmd) > 0 {
			backgroundMsg += "`" + prefix + bg.Line[i].Cmd + "`"
		}
		backgroundMsg += bg.Line[i].Post + "\n"
	}
}
