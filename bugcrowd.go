package main

import (
	"encoding/json"
	"fmt"
)

const BugcrowdURL = "https://raw.githubusercontent.com/Osb0rn3/bugbounty-targets/main/programs/bugcrowd.json"

type BugcrowdList []struct {
	Name         string `json:"name"`
	ProgramURL   string `json:"program_url"`
	LicenseKey   string `json:"license_key"`
	TargetGroups []struct {
		InScope bool `json:"in_scope"`
		Targets []struct {
			Name     string `json:"name"`
			Category string `json:"category"`
			URI      string `json:"uri"`
			Target   struct {
				Tags []struct {
					Name string `json:"name"`
				} `json:"tags"`
			} `json:"target"`
		} `json:"targets"`
	} `json:"target_groups"`
	MinRewards string `json:"min_rewards,omitempty"`
	MaxRewards string `json:"max_rewards,omitempty"`
}

func GetBugcrowdChanges(file string) []companyChanges {
	changes := []companyChanges{}
	rawBody := getReq(BugcrowdURL)
	// check if the file already exists or not
	if fileExists(file) {
		// unmarshal requested body to the related struct
		var newData BugcrowdList
		err := json.Unmarshal(rawBody, &newData)
		checkError(err)

		// unmarshal old file data to the related struct
		var oldData BugcrowdList
		content := loadFile(file)
		err = json.Unmarshal(content, &oldData)
		checkError(err)

		//use all old targets' name + uri as map index
		allOldTargets := make(map[string]bool)
		for _, oldProgram := range oldData {
			for _, oldTargetGroups := range oldProgram.TargetGroups {
				if oldTargetGroups.InScope {
					for _, oldTarget := range oldTargetGroups.Targets {
						allOldTargets[oldTarget.Name+oldTarget.URI] = true
					}
				}
			}
		}

		// check if new targets are used as a key in the allOldTargets map
		for _, newProgram := range newData {
			var assetChanges []string

			//check if each target is present in the old database
			for _, newTargetGroups := range newProgram.TargetGroups {
				if newTargetGroups.InScope {
					for _, newTarget := range newTargetGroups.Targets {
						if _, exist := allOldTargets[newTarget.Name+newTarget.URI]; !exist {
							assetChanges = append(assetChanges, fmt.Sprintf("%s\t%s", newTarget.Name, newTarget.URI))
						}
					}
				}
			}

			// add changes if a company has new assets
			if len(assetChanges) != 0 {
				var programType string
				if newProgram.LicenseKey == "vdp_pro" || newProgram.LicenseKey == "vdp_essentials" {
					programType = "vdp"
				} else {
					programType = "bbp"
				}
				changes = append(changes, companyChanges{newProgram.Name, programType, fmt.Sprintf("https://bugcrowd.com%s", newProgram.ProgramURL), assetChanges})
			}
		}

		// replace the old data with the new one (the one received from a get request)
		saveToFile(file, rawBody)
	} else {
		saveToFile(file, rawBody)
	}

	return changes
}
