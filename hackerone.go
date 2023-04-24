package main

import (
	"encoding/json"
	"fmt"
	"time"
)

const HackeroneURL = "https://raw.githubusercontent.com/Osb0rn3/bugbounty-targets/main/programs/hackerone.json"

type HackeroneList []struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Handle          string `json:"handle"`
		Name            string `json:"name"`
		SubmissionState string `json:"submission_state"`
		State           string `json:"state"`
		OffersBounties  bool   `json:"offers_bounties"`
	} `json:"attributes"`
	Relationships struct {
		StructuredScopes struct {
			Data []struct {
				ID         string `json:"id"`
				Type       string `json:"type"`
				Attributes struct {
					AssetType             string    `json:"asset_type"`
					AssetIdentifier       string    `json:"asset_identifier"`
					EligibleForBounty     bool      `json:"eligible_for_bounty"`
					EligibleForSubmission bool      `json:"eligible_for_submission"`
					Instruction           string    `json:"instruction"`
					MaxSeverity           string    `json:"max_severity"`
					CreatedAt             time.Time `json:"created_at"`
					UpdatedAt             time.Time `json:"updated_at"`
				} `json:"attributes,omitempty"`
			} `json:"data"`
		} `json:"structured_scopes"`
	} `json:"relationships"`
}

func GetHackeroneChanges(file string) []companyChanges {
	changes := []companyChanges{}
	rawBody := getReq(HackeroneURL)
	// check if the file already exists or not
	if fileExists(file) {
		// unmarshal requested body to the related struct
		var newData HackeroneList
		err := json.Unmarshal(rawBody, &newData)
		checkError(err)

		// unmarshal old file data to the related struct
		var oldData HackeroneList
		content := loadFile(file)
		err = json.Unmarshal(content, &oldData)
		checkError(err)

		//use all old targets' asset identifier as map index
		allOldTargets := make(map[string]bool)
		for _, oldProgram := range oldData {
			for _, oldTarget := range oldProgram.Relationships.StructuredScopes.Data {
				if oldTarget.Attributes.EligibleForSubmission {
					allOldTargets[oldTarget.Attributes.AssetIdentifier] = true
				}
			}
		}

		// check if new targets are used as a key in the allOldTargets map
		for _, newProgram := range newData {
			var assetChanges []string

			//check if each target is present in the old database
			for _, newTarget := range newProgram.Relationships.StructuredScopes.Data {
				if newTarget.Attributes.EligibleForSubmission {
					if _, exist := allOldTargets[newTarget.Attributes.AssetIdentifier]; !exist {
						assetChanges = append(assetChanges, newTarget.Attributes.AssetIdentifier)
					}
				}
			}

			// add changes if a company has new assets
			if len(assetChanges) != 0 {
				changes = append(changes, companyChanges{newProgram.Attributes.Name, fmt.Sprintf("https://hackerone.com/%s", newProgram.Attributes.Handle), assetChanges})
			}
		}

		// replace the old data with the new one (the one received from a get request)
		saveToFile(file, rawBody)
	} else {
		saveToFile(file, rawBody)
	}

	return changes
}
