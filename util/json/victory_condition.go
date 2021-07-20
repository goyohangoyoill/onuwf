// Package json is a package for json files
package json

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	vcTitle string
	vcMsg   string
)

type victoryCondition struct {
	Title string `json:"title"`
	Team  []team `json:"team"`
}

type team struct {
	Title string   `json:"title"`
	Line  []string `json:"line"`
}

// victory_condition.json 파일 읽어서 "ㅁ승리조건" 실행시 출력할 데이터 세팅
func readVictoryConditionJSON() {
	jsonFile, err := os.Open("./asset/victory_condition.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	var vc victoryCondition
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(byteValue, &vc)
	vcTitle = "**" + vc.Title + "**"
	vcMsg = ""
	for _, team := range vc.Team {
		vcMsg += "**" + team.Title + "**" + "\n"
		for _, line := range team.Line {
			if len(line) > 0 {
				vcMsg += "> " + line + "\n"
			}
		}
		vcMsg += "\n"
	}
}
