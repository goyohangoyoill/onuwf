/* "ㅁ승리조건" 명령어 관련 함수 */
package util

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
	jsonFile, err := os.Open("./asset/victory_condition.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	var vc victory_condition
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(byteValue, &vc)

	vcTitle = "**" + vc.Title + "**"
	vcMsg = ""
	for i := 0; i < len(vc.Team); i++ {
		vcMsg += "**" + vc.Team[i].Title + "**" + "\n"
		for j := 0; j < len(vc.Team[i].Line); j++ {
			vcMsg += vc.Team[i].Line[j] + "\n"
		}
		vcMsg += "\n"
	}
}
