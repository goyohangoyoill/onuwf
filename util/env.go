// Package util is a package for json files and database.
package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// EnvInit 설치 환경 불러오기.
func EnvInit() map[string]string {
	envFile, err := os.Open("asset/env.json")
	if err != nil {
		log.Fatal(err)
	}
	defer envFile.Close()
	var byteValue []byte
	byteValue, err = ioutil.ReadAll(envFile)
	if err != nil {
		log.Fatal(err)
	}
	env := make(map[string]string)
	json.Unmarshal([]byte(byteValue), &env)
	return env
}
