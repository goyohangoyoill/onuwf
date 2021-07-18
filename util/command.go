// Package util is a package for json files and database.
package util

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
	for i := 0; i < len(cmd.List); i++ {
		commandMsg += "`" + prefix + cmd.List[i].Name + "` : "
		for j := 0; j < len(cmd.List[i].Guide); j++ {
			commandMsg += cmd.List[i].Guide[j] + "\n"
		}
		if i == 4 || i == 7 || i == 10 {
			commandMsg += "\n"
		}
	}
}
