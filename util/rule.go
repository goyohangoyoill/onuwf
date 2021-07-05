/* "ㅁ게임방법" 명령어 관련 함수 */
package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	ruleTitle string
	ruleMsg   string
)

type rule struct {
	Title string `json:"title"`
	Line  []line `json:"line"`
}

type line struct {
	Pre  string `json:"pre"`
	Cmd  string `json:"cmd"`
	Post string `json:"post"`
}

// rule.json 파일 읽어서 "ㅁ게임방법" 실행시 출력할 데이터 세팅
func readRuleJSON() {
	jsonFile, err := os.Open("./asset/rule.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer jsonFile.Close()
	var rule rule
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(byteValue, &rule)

	ruleTitle = "**" + rule.Title + "**"
	ruleMsg = ""
	for i := 0; i < len(rule.Line); i++ {
		if len(rule.Line[i].Pre) > 0 {
			ruleMsg += rule.Line[i].Pre
		}
		if len(rule.Line[i].Cmd) > 0 {
			ruleMsg += "`" + prefix + rule.Line[i].Cmd + "`"
		}
		ruleMsg += rule.Line[i].Post + "\n"
	}
}
