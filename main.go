package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	discordURL_BBP string = "BBP WEBHOOK URL"
	discordURL_VDP string = "VDP WEBHOOK URL"
)

type companyChanges struct {
	Name    string
	Type    string
	URL     string
	Changes []string
}

type webhookPayload struct {
	Contents string `json:"content"`
}

func main() {

	// create a directory for json files
	createDirIfNotExist("platforms")

	// compare bugcrowd changes
	bugcrowdChanges := GetBugcrowdChanges("platforms/bugcrowd.json")
	// send changes to the discord server
	sendNotif(bugcrowdChanges)

	// compare hackerone changes
	hackeroneChanges := GetHackeroneChanges("platforms/hackerone.json")
	// send changes to the discord server
	sendNotif(hackeroneChanges)

	// compare yeswehack changes
	yeswehackChanges := GetYeswehackChanges("platforms/yeswehack.json")
	// send changes to the discord server
	sendNotif(yeswehackChanges)

	// compare intigriti changes
	intigritiChanges := GetIntigritiChanges("platforms/intigriti.json")
	// send changes to the discord server
	sendNotif(intigritiChanges)

	//log this program's activity
	healthCheck()
}

func createDirIfNotExist(path string) {
	// check if directory already exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// if directory doesn't exist, create it
		err = os.MkdirAll(path, 0755)
		checkError(err)
	}
}

func sendNotif(changes []companyChanges) {
	//beautify it for discord markdown
	for _, company := range changes {
		content := fmt.Sprintf(">>> **%v**\n**Type**: %v\n**URL**: <%v>\n**Assets**:\n```\n", company.Name, company.Type, company.URL)
		for _, asset := range company.Changes {
			content += asset + "\n"
		}
		content += "```"

		postData, err := json.Marshal(webhookPayload{Contents: content})
		checkError(err)

		if company.Type == "vdp" {
			http.Post(discordURL_VDP, "application/json", bytes.NewBuffer(postData))
		} else {
			http.Post(discordURL_BBP, "application/json", bytes.NewBuffer(postData))
		}
	}
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
		log.Panic(err)
	}
}

func getReq(url string) []byte {
	response, err := http.Get(url)
	checkError(err)
	defer response.Body.Close()
	rawBody, err := ioutil.ReadAll(response.Body)
	checkError(err)

	return rawBody
}

func saveToFile(path string, content []byte) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	checkError(err)
	defer file.Close()
	io.Copy(file, bytes.NewBuffer(content))
}

func loadFile(path string) []byte {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	checkError(err)
	content, err := io.ReadAll(file)
	checkError(err)

	return content
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
