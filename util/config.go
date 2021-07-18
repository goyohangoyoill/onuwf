package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// ONUWF/asset/config.json 파일 읽어서 return하는 함수
func ReadConfigJson() map[string]string {
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
	config := make(map[string]string)
	json.Unmarshal([]byte(byteValue), &config)
	return config
}
