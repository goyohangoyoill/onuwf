package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	Prefix       string `json:"prefix"`
	VoteDelaySec int    `json:"voteDelaySec"`
}

// ONUWF/asset/config.json 파일 읽어서 return하는 함수
func ReadConfigJson() Config {
	configFile, err := os.Open("asset/config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer configFile.Close()

	var byteValue []byte
	byteValue, err = ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal(err)
	}
	var config Config
	json.Unmarshal([]byte(byteValue), &config)

	// 음수 시간 동안 기다릴 순 없으니
	if config.VoteDelaySec < 0 {
		config.VoteDelaySec = 0
	}
	return config
}
