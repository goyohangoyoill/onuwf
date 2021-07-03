/* "ㅁ참고" 명령어 관련 함수 */
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	Note      note
	noteTitle string
	noteMsg   string
)

type note struct {
	Title string     `json:"title"`
	Line  []noteLine `json:"line"`
}

type noteLine struct {
	Bold string `json:"bold"`
	Post string `json:"post"`
}

// note.json 파일 읽어서 "ㅁ참고" 실행시 출력할 데이터 세팅
func readNoteJSON() {
	noteFile, err := os.Open("asset/note.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer noteFile.Close()
	byteValue, err := ioutil.ReadAll(noteFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(byteValue, &Note)

	noteTitle = "**" + Note.Title + "**"
	noteMsg = ""
	for i := 0; i < len(Note.Line); i++ {
		if len(Note.Line[i].Bold) > 0 {
			noteMsg += "**" + Note.Line[i].Bold + "**"
		}
		noteMsg += Note.Line[i].Post + "\n"
	}
	list := roleList()
	for i, item := range list {
		noteMsg += item + " "
		if i%4 == 3 {
			noteMsg += "\n"
		}
	}
}
