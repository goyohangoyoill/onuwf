/* "ㅁ승리조건" 명령어 관련 함수 */
package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	VC      victory_condition
	VCTitle string
	VCMsg   string
)

type victory_condition struct {
	Title string `json:"title"`
	Team  []team `json:"team"`
}

type team struct {
	Title string   `json:"title"`
	Line  []string `json:"line"`
}

// victory_condition.json 파일 읽어서 "ㅁ승리조건" 실행시 출력할 데이터 세팅
func readVictoryConditionJSON() {
	jsonFile, err := os.Open("asset/victory_condition.json")
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
	json.Unmarshal(byteValue, &VC)

	VCTitle = "**" + VC.Title + "**"
	VCMsg = ""
	for i := 0; i < len(VC.Team); i++ {
		VCMsg += "**" + VC.Team[i].Title + "**" + "\n"
		for j := 0; j < len(VC.Team[i].Line); j++ {
			VCMsg += VC.Team[i].Line[j] + "\n"
		}
		VCMsg += "\n"
	}
}
