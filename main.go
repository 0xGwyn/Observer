package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

const discordURL string = "DISCORD WEBHOOK URL"

func main() {

	platformUrls := map[string]string{
		"hackerone": "https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/hackerone_data.json",
		"bugcrowd":  "https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/bugcrowd_data.json",
		"intigriti": "https://raw.githubusercontent.com/arkadiyt/bounty-targets-data/main/data/intigriti_data.json",
	}

	for platformName, url := range platformUrls {
		body := getReq(url)

		//either create the file if does not exist or check for changes
		if fileExists(platformName + ".json") {
			content := loadFileToString(platformName + ".json")

			//compare the most recent changes with the old one
			changes := compareTargets(content, body, platformName)

			//send changes to the discord server
			sendNotif(discordURL, platformName, changes)

			//replace newest data with the old one
			saveStringToFile(platformName+".json", body)

		} else {
			saveStringToFile(platformName+".json", body)
		}

	}
}

type companyChanges struct {
	name    string
	url     string
	changes []string
}

func sendNotif(url, platform string, changes []companyChanges) {
	// to be implemented
}

func compareTargets(old, new, platform string) []companyChanges {
	changes := []companyChanges{}

	//remove junk characters for comparison
	old = strings.Replace(strings.Replace(strings.Replace(old, "\\t", "", -1), "\\\"", "", -1), "", "", -1)

	switch platform {

	case "bugcrowd":
		newData := stringToStruct[Bugcrowd](new)

		for _, company := range newData {
			var assetChanges []string

			for _, inscope := range company.Targets.InScope {
				//remove junk characters for comparison
				target := strings.Replace(strings.Replace(strings.Replace(inscope.Target, "\t", "", -1), "\"", "", -1), "", "", -1)

				if !strings.Contains(old, target) {
					assetChanges = append(assetChanges, target)
				}
			}

			// add changes if a company has new assets
			if len(assetChanges) != 0 {
				changes = append(changes, companyChanges{company.Name, company.URL, assetChanges})
			}
		}

	case "hackerone":
		newData := stringToStruct[Hackerone](new)

		for _, company := range newData {
			var assetChanges []string

			for _, inscope := range company.Targets.InScope {
				//remove junk characters for comparison
				target := strings.Replace(strings.Replace(strings.Replace(inscope.AssetIdentifier, "\t", "", -1), "\"", "", -1), "", "", -1)

				if !strings.Contains(old, target) {
					assetChanges = append(assetChanges, target)
				}
			}

			// add changes if a company has new assets
			if len(assetChanges) != 0 {
				changes = append(changes, companyChanges{company.Name, company.URL, assetChanges})
			}
		}

	case "intigriti":
		newData := stringToStruct[Intigriti](new)

		for _, company := range newData {
			var assetChanges []string

			for _, inscope := range company.Targets.InScope {
				//remove junk characters for comparison
				target := strings.Replace(strings.Replace(strings.Replace(inscope.Endpoint, "\t", "", -1), "\"", "", -1), "", "", -1)

				if !strings.Contains(old, target) {
					assetChanges = append(assetChanges, target)
				}
			}

			// add changes if a company has new assets
			if len(assetChanges) != 0 {
				changes = append(changes, companyChanges{company.Name, company.URL, assetChanges})
			}
		}
	}

	return changes
}

type platformtype interface {
	Bugcrowd | Hackerone | Intigriti
}

func stringToStruct[platform platformtype](data string) []platform {
	var p []platform
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
	// err := ioutil.WriteFile(path, []byte(content), 0644)
	checkError(err)
	defer file.Close()
	fmt.Fprint(file, content)
}

func loadFileToString(path string) string {
	content, err := ioutil.ReadFile(path)
	checkError(err)
	return string(content)
}
