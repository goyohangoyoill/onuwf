// Package util is a package for json files and database.
package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	noteTitle string
	noteMsg   string
)

type note struct {
	Title string     `json:"title"`
	Line  []noteLine `json:"line"`
}

type noteLine struct {
	Bold bool   `json:"bold"`
	Msg  string `json:"msg"`
}

// note.json 파일 읽어서 "ㅁ참고" 실행시 출력할 데이터 세팅
func readNoteJSON(rg []RoleGuide) {
	noteFile, err := os.Open("./asset/note.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer noteFile.Close()
	var note note
	byteValue, err := ioutil.ReadAll(noteFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(byteValue, &note)

	noteTitle = "**" + note.Title + "**"
	noteMsg = ""
	for i := 0; i < len(note.Line); i++ {
		switch note.Line[i].Bold {
		case true:
			noteMsg += "**" + note.Line[i].Msg + "**"
		case false:
			noteMsg += note.Line[i].Msg
		}
		noteMsg += "\n"
	}
	list := roleList(rg)
	for i, item := range list {
		noteMsg += item + " "
		if i%4 == 3 {
			noteMsg += "\n"
		}
	}
}
