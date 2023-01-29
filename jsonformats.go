package main

type Bugcrowd struct {
	Name              string `json:"name"`
	URL               string `json:"url"`
	AllowsDisclosure  bool   `json:"allows_disclosure"`
	ManagedByBugcrowd bool   `json:"managed_by_bugcrowd"`
	SafeHarbor        string `json:"safe_harbor"`
	MaxPayout         int    `json:"max_payout"`
	Targets           struct {
		InScope []struct {
			Type   string `json:"type"`
			Target string `json:"target"`
		} `json:"in_scope"`
		OutOfScope []struct {
			Type   string `json:"type"`
			Target string `json:"target"`
		} `json:"out_of_scope"`
	} `json:"targets"`
}

type Intigriti struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	CompanyHandle        string `json:"company_handle"`
	Handle               string `json:"handle"`
	URL                  string `json:"url"`
	Status               string `json:"status"`
	ConfidentialityLevel string `json:"confidentiality_level"`
	MinBounty            struct {
		Value    float64 `json:"value"`
		Currency string  `json:"currency"`
	} `json:"min_bounty"`
	MaxBounty struct {
		Value    float64 `json:"value"`
		Currency string  `json:"currency"`
	} `json:"max_bounty"`
	Targets struct {
		InScope []struct {
			Type        string `json:"type"`
			Endpoint    string `json:"endpoint"`
			Description string `json:"description"`
		} `json:"in_scope"`
		OutOfScope []struct {
			Type        string `json:"type"`
			Endpoint    string `json:"endpoint"`
			Description string `json:"description"`
		} `json:"out_of_scope"`
	} `json:"targets"`
}

type Hackerone struct {
	AllowsBountySplitting             bool    `json:"allows_bounty_splitting"`
	AverageTimeToBountyAwarded        float64 `json:"average_time_to_bounty_awarded"`
	AverageTimeToFirstProgramResponse float64 `json:"average_time_to_first_program_response"`
	AverageTimeToReportResolved       float64 `json:"average_time_to_report_resolved"`
	Handle                            string  `json:"handle"`
	ID                                int     `json:"id"`
	ManagedProgram                    bool    `json:"managed_program"`
	Name                              string  `json:"name"`
	OffersBounties                    bool    `json:"offers_bounties"`
	OffersSwag                        bool    `json:"offers_swag"`
	ResponseEfficiencyPercentage      int     `json:"response_efficiency_percentage"`
	SubmissionState                   string  `json:"submission_state"`
	URL                               string  `json:"url"`
	Website                           string  `json:"website"`
	Targets                           struct {
		InScope []struct {
			AssetIdentifier            string `json:"asset_identifier"`
			AssetType                  string `json:"asset_type"`
			AvailabilityRequirement    string `json:"availability_requirement"`
			ConfidentialityRequirement string `json:"confidentiality_requirement"`
			EligibleForBounty          bool   `json:"eligible_for_bounty"`
			EligibleForSubmission      bool   `json:"eligible_for_submission"`
			Instruction                string `json:"instruction"`
			IntegrityRequirement       string `json:"integrity_requirement"`
			MaxSeverity                string `json:"max_severity"`
		} `json:"in_scope"`
		OutOfScope []struct {
			AssetIdentifier            string `json:"asset_identifier"`
			AssetType                  string `json:"asset_type"`
			AvailabilityRequirement    string `json:"availability_requirement"`
			ConfidentialityRequirement string `json:"confidentiality_requirement"`
			EligibleForBounty          bool   `json:"eligible_for_bounty"`
			EligibleForSubmission      bool   `json:"eligible_for_submission"`
			Instruction                string `json:"instruction"`
			IntegrityRequirement       string `json:"integrity_requirement"`
			MaxSeverity                string `json:"max_severity"`
		} `json:"out_of_scope"`
	} `json:"targets"`
}
