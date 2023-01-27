package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {

	platformUrls := map[string]string{
		"hackerone": "https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/hackerone_data.json",
		"bugcrowd":  "https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/bugcrowd_data.json",
		"intigriti": "https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/intigriti_data.json",
	}
	for platformName, url := range platformUrls {
		body := getReq(url)
		if fileExists(platformName + ".json") {
			//here we have to compare two jsons, send diff notif then replace the new struct/json with .json file
			// content := loadFile(platformName + ".json")
			// oldData, err := json.Marshal(content)
			// checkError(err)

			// var newData interface{}
			// switch platformName {
			// case "bugcrowd":
			// 	newData = stringToStruct[Bugcrowd](body)
			// case "hackerone":
			// 	newData = stringToStruct[Hackerone](body)
			// case "intigriti":
			// 	newData = stringToStruct[Intigriti](body)
			// }

		} else {
			saveStringToFile(platformName+".json", body)
		}
	}
}

func stringToStruct[platforms Bugcrowd | Hackerone | Intigriti](data string) platforms {
	var p platforms
	err := json.Unmarshal([]byte(data), &p)
	checkError(err)
	return p
}

func fileExists(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true
	} else {
		return false
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func getReq(url string) string {
	response, err := http.Get(url)
	checkError(err)
	defer response.Body.Close()
	rawBody, err := ioutil.ReadAll(response.Body)
	checkError(err)
	return string(rawBody)
}

func saveStringToFile(path string, content string) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0644)
	checkError(err)
	defer file.Close()
	fmt.Fprint(file, content)
}

func loadFile(path string) string {
	content, err := ioutil.ReadFile(path)
	checkError(err)
	return string(content)
}
