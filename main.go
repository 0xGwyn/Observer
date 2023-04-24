package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	discordURL string = "DISCORD'S WEBHOOK URL"
	chaosUrl          = "https://raw.githubusercontent.com/projectdiscovery/public-bugbounty-programs/main/chaos-bugbounty-list.json"
)

// ChaosProgram json data item struct
type ChaosProgram struct {
	Name    string   `json:"name"`
	URL     string   `json:"url"`
	Bounty  bool     `json:"bounty"`
	Swag    bool     `json:"swag"`
	Domains []string `json:"domains"`
}

type ChaosList struct {
	Programs []ChaosProgram `json:"programs"`
}

type companyChanges struct {
	name    string
	url     string
	changes []string
}

type webhookPayload struct {
	Contents string `json:"content"`
}

func main() {

	//request json data from github
	resp := getReq(chaosUrl)

	//either create the file if does not exist or check for changes
	if fileExists("chaos.json") {
		content := loadFileToString("chaos.json")

		//compare the most recent changes with the old one
		changes := compareTargets(content, resp)

		//send changes to the discord server
		sendNotif(changes)

		//replace new data with the old one
		saveStringToFile("chaos.json", resp)

	} else {
		saveStringToFile("chaos.json", resp)
	}

	//log program's activity
	healthCheck()
}

func sendNotif(changes []companyChanges) {
	//beautify it for discord markup
	for _, company := range changes {
		content := fmt.Sprintf("**%v**\n*URL*: <%v>\n*Assets*:\n```\n", company.name, company.url)
		for _, asset := range company.changes {
			content = content + asset + "\n"
		}
		content += "```"

		postData, err := json.Marshal(webhookPayload{Contents: content})
		checkError(err)

		http.Post(discordURL, "application/json", bytes.NewBuffer(postData))
	}
}

func compareTargets(old, new string) []companyChanges {
	changes := []companyChanges{}

	//unmarshal old and new data as ChaosList
	newData := stringToStruct(new)
	oldData := stringToStruct(old)

	//use all old domains as map index
	allOldDomains := make(map[string]bool)
	for _, oldProgram := range oldData.Programs {
		for _, oldDomain := range oldProgram.Domains {
			allOldDomains[oldDomain] = true
		}
	}

	for _, newProgram := range newData.Programs {
		var assetChanges []string

		//check if each domain is present in the old database
		for _, newDomain := range newProgram.Domains {
			if _, exist := allOldDomains[newDomain]; !exist {
				assetChanges = append(assetChanges, newDomain)
			}
		}

		// add changes if a company has new assets
		if len(assetChanges) != 0 {
			changes = append(changes, companyChanges{newProgram.Name, newProgram.URL, assetChanges})
		}
	}

	return changes
}

func stringToStruct(data string) ChaosList {
	var s ChaosList
	err := json.Unmarshal([]byte(data), &s)
	checkError(err)
	return s
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

func saveStringToFile(path, content string) {
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

func healthCheck() {

	file, err := os.OpenFile("health_check", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	time := getCurrentTime()
	file.WriteString(time + " working fine\n")
}

func getCurrentTime() string {
	currentTime := time.Now()
	formattedTime := fmt.Sprintf("%v/%v/%v %v:%v:%v", currentTime.Year(), int(currentTime.Month()), currentTime.Day(), currentTime.Hour(), currentTime.Minute(), currentTime.Second())

	return formattedTime
}
