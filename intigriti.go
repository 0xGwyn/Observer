package main

import (
	"encoding/json"
	"fmt"
)

const IntigritiURL = "https://raw.githubusercontent.com/Osb0rn3/bugbounty-targets/main/programs/intigriti.json"

type IntigritiList []struct {
	ProgramID          string `json:"programId"`
	CompanyHandle      string `json:"companyHandle"`
	CompanyName        string `json:"companyName"`
	TacRequired        bool   `json:"tacRequired"`
	CompanySustainable bool   `json:"companySustainable"`
	Handle             string `json:"handle"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	MinBounty          struct {
		Value    float64 `json:"value"`
		Currency string  `json:"currency"`
	} `json:"minBounty"`
	MaxBounty struct {
		Value    float64 `json:"value"`
		Currency string  `json:"currency"`
	} `json:"maxBounty"`
	LastUpdatedAt    int `json:"lastUpdatedAt"`
	LastSubmissionAt int `json:"lastSubmissionAt"`
	Domains          []struct {
		Type           int    `json:"type"`
		Endpoint       string `json:"endpoint"`
		BusinessImpact int    `json:"businessImpact"`
		Description    string `json:"description"`
	} `json:"domains"`
}

func GetIntigritiChanges(file string) []companyChanges {
	changes := []companyChanges{}
	rawBody := getReq(IntigritiURL)
	// check if the file already exists or not
	if fileExists(file) {
		// unmarshal requested body to the related struct
		var newData IntigritiList
		err := json.Unmarshal(rawBody, &newData)
		checkError(err)

		// unmarshal old file data to the related struct
		var oldData IntigritiList
		content := loadFile(file)
		err = json.Unmarshal(content, &oldData)
		checkError(err)

		//use all old targets' endpoints as map index
		allOldTargets := make(map[string]bool)
		for _, oldProgram := range oldData {
			for _, oldTarget := range oldProgram.Domains {
				allOldTargets[oldTarget.Endpoint] = true
			}
		}

		// check if new targets are used as a key in the allOldTargets map
		for _, newProgram := range newData {
			var assetChanges []string

			//check if each target is present in the old database
			for _, newTarget := range newProgram.Domains {
				if _, exist := allOldTargets[newTarget.Endpoint]; !exist {
					assetChanges = append(assetChanges, newTarget.Endpoint)
				}
			}

			// add changes if a company has new assets
			if len(assetChanges) != 0 {
				var programType string
				if newProgram.MaxBounty.Value == 0 && newProgram.MinBounty.Value == 0 {
					programType = "vdp"
				} else {
					programType = "bug_bounty"
				}
				changes = append(changes, companyChanges{newProgram.Name, programType, fmt.Sprintf("https://app.intigriti.com/programs/%s/%s/detail", newProgram.CompanyHandle, newProgram.Handle), assetChanges})
			}
		}

		// replace the old data with the new one (the one received from a get request)
		saveToFile(file, rawBody)
	} else {
		saveToFile(file, rawBody)
	}

	return changes
}
