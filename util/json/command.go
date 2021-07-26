// Package json is a package for json files
package json

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	commandTitle string
	commandMsg   string
)

type command struct {
	Title string        `json:"title"`
	List  []commandList `json:"list"`
}

type commandList struct {
	Name  string   `json:"name"`
	Guide []string `json:"guide"`
}

// command.json 파일 읽어서 "ㅁ명령어" 실행시 출력할 데이터 세팅
func readCommandJSON(prefix string) {
	cmdFile, err := os.Open("./asset/command.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cmdFile.Close()
	var cmd command
	byteValue, err := ioutil.ReadAll(cmdFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(byteValue, &cmd)
	commandTitle = "**" + cmd.Title + "**"
	commandMsg = ""
	for _, list := range cmd.List {
		if len(list.Name) > 0 {
			commandMsg += "`" + prefix + list.Name + "` : "
		}
		for _, guide := range list.Guide {
			commandMsg += guide + "\n"
		}
	}
}
