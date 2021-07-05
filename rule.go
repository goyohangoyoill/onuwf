/* "ㅁ게임방법" 명령어 관련 함수 */
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	ruleTitle string
	ruleMsg   string
	Rule      rule
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
	jsonFile, err := os.Open("asset/rule.json")
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
	json.Unmarshal(byteValue, &Rule)

	ruleTitle = "**" + Rule.Title + "**"
	ruleMsg = ""
	for i := 0; i < len(Rule.Line); i++ {
		if len(Rule.Line[i].Pre) > 0 {
			ruleMsg += Rule.Line[i].Pre
		}
		if len(Rule.Line[i].Cmd) > 0 {
			ruleMsg += "`" + prefix + Rule.Line[i].Cmd + "`"
		}
		ruleMsg += Rule.Line[i].Post + "\n"
	}
}
