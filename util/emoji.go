package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// EmojiInit 이모지 맵에 불러오기.
func EmojiInit() map[string]string {
	emjFile, err := os.Open("asset/emoji.json")
	if err != nil {
		log.Fatal(err)
	}
	defer emjFile.Close()

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(emjFile)
	if err != nil {
		log.Fatal(err)
	}
	emj := make(map[string]string)
	json.Unmarshal([]byte(byteValue), &emj)
	return emj
}
