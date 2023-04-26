package main

import (
	"encoding/json"
	"fmt"
)

const YeswehackURL = "https://raw.githubusercontent.com/Osb0rn3/bugbounty-targets/main/programs/yeswehack.json"

type YeswehackList []struct {
	Title           string `json:"title"`
	Slug            string `json:"slug"`
	Public          bool   `json:"public"`
	Bounty          bool   `json:"bounty"`
	Gift            bool   `json:"gift"`
	Disabled        bool   `json:"disabled"`
	Vdp             bool   `json:"vdp"`
	Archived        bool   `json:"archived"`
	BountyRewardMin int    `json:"bounty_reward_min"`
	BountyRewardMax int    `json:"bounty_reward_max"`
	ScopesCount     int    `json:"scopes_count"`
	Scopes          []struct {
		Scope               string `json:"scope"`
		ScopeType           string `json:"scope_type"`
		SecurityRequirement string `json:"security_requirement"`
	} `json:"scopes"`
}

func GetYeswehackChanges(file string) []companyChanges {
	changes := []companyChanges{}
	rawBody := getReq(YeswehackURL)
	// check if the file already exists or not
	if fileExists(file) {
		// unmarshal requested body to the related struct
		var newData YeswehackList
		err := json.Unmarshal(rawBody, &newData)
		checkError(err)

		// unmarshal old file data to the related struct
		var oldData YeswehackList
		content := loadFile(file)
		err = json.Unmarshal(content, &oldData)
		checkError(err)

		//use all old targets' scopes as map index
		allOldTargets := make(map[string]bool)
		for _, oldProgram := range oldData {
			for _, oldTarget := range oldProgram.Scopes {
				allOldTargets[oldTarget.Scope] = true
			}
		}

		// check if new targets are used as a key in the allOldTargets map
		for _, newProgram := range newData {
			var assetChanges []string

			//check if each target is present in the old database
			for _, newTarget := range newProgram.Scopes {
				if _, exist := allOldTargets[newTarget.Scope]; !exist {
					assetChanges = append(assetChanges, newTarget.Scope)
				}
			}

			// add changes if a company has new assets
			if len(assetChanges) != 0 {
				var programType string
				if newProgram.Vdp {
					programType = "vdp"
				} else {
					programType = "bug_bounty"
				}
				changes = append(changes, companyChanges{newProgram.Title, programType, fmt.Sprintf("https://yeswehack.com/programs/%s", newProgram.Slug), assetChanges})
			}
		}

		// replace the old data with the new one (the one received from a get request)
		saveToFile(file, rawBody)
	} else {
		saveToFile(file, rawBody)
	}

	return changes
}
