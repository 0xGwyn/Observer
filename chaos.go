package main

import (
	"encoding/json"
)

const ChaosURL = "https://raw.githubusercontent.com/projectdiscovery/public-bugbounty-programs/main/chaos-bugbounty-list.json"

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

func GetChaosChanges(file string) []companyChanges {
	changes := []companyChanges{}
	rawBody := getReq(ChaosURL)
	// check if the file already exists or not
	if fileExists(file) {
		// unmarshal requested body to the related struct
		var newData ChaosList
		err := json.Unmarshal(rawBody, &newData)
		checkError(err)

		// unmarshal old file data to the related struct
		var oldData ChaosList
		content := loadFile(file)
		err = json.Unmarshal(content, &oldData)
		checkError(err)

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
				var programType string
				if newProgram.Bounty {
					programType = "bug_bounty"
				} else {
					programType = "vdp"
				}
				changes = append(changes, companyChanges{newProgram.Name, programType, newProgram.URL, assetChanges})
			}
		}

		// replace the old data with the new one (the one received from a get request)
		saveToFile(file, rawBody)
	} else {
		saveToFile(file, rawBody)
	}

	return changes
}
