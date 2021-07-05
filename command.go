/* "ㅁ명령어" 명령어 관련 함수 */
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	commandTitle string
	commandMsg   string
	Command      command
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
func readCommandJSON() {
	cmdFile, err := os.Open("asset/command.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cmdFile.Close()
	byteValue, err := ioutil.ReadAll(cmdFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(byteValue, &Command)

	commandTitle = "**" + Command.Title + "**"
	commandMsg = ""
	for i := 0; i < len(Command.List); i++ {
		commandMsg += "`" + prefix + Command.List[i].Name + "` : "
		for j := 0; j < len(Command.List[i].Guide); j++ {
			commandMsg += Command.List[i].Guide[j] + "\n"
		}
		if i == 4 || i == 7 || i == 10 {
			commandMsg += "\n"
		}
	}
}
