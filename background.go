/* "ㅁ게임배경" 명령어 관련 함수 */
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	Background      background
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
func readBackgroundJSON() {
	jsonFile, err := os.Open("asset/background.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(byteValue, &Background)

	backgroundTitle = "**" + Background.Title + "**"
	backgroundMsg = ""
	for i := 0; i < len(Background.Line); i++ {
		if len(Background.Line[i].Pre) > 0 {
			backgroundMsg += Background.Line[i].Pre
		}
		if len(Background.Line[i].Cmd) > 0 {
			backgroundMsg += "`" + prefix + Background.Line[i].Cmd + "`"
		}
		backgroundMsg += Background.Line[i].Post + "\n"
	}
}
