package diagnostic

func SARIF(toolName string, items []Item) sarifReport {
	report := sarifReport{
		Version: "2.1.0",
		Schema:  "https://json.schemastore.org/sarif-2.1.0.json",
		Runs: []sarifRun{{
			Tool: sarifTool{
				Driver: sarifDriver{
					Name: toolName,
				},
			},
		}},
	}

	seenRules := map[string]bool{}
	for _, item := range items {
		if !seenRules[item.RuleID] {
			seenRules[item.RuleID] = true
			rule := sarifRule{ID: item.RuleID}
			rule.ShortDescription.Text = item.Message
			rule.Help.Text = item.FixHint
			report.Runs[0].Tool.Driver.Rules = append(report.Runs[0].Tool.Driver.Rules, rule)
		}

		result := sarifResult{RuleID: item.RuleID, Level: string(item.Severity)}
		result.Message.Text = item.Message
		if item.File != "" {
			result.Locations = []sarifLocation{{
				PhysicalLocation: sarifPhysicalLocation{
					ArtifactLocation: sarifArtifactLocation{URI: item.File},
					Region:           sarifRegion{StartLine: max(item.Line, 1)},
				},
			}}
		}
		report.Runs[0].Results = append(report.Runs[0].Results, result)
	}

	return report
}

func max(left, right int) int {
	if left > right {
		return left
	}
	return right
}

type sarifReport struct {
	Version string     `json:"version"`
	Schema  string     `json:"$schema"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name  string      `json:"name"`
	Rules []sarifRule `json:"rules"`
}

type sarifRule struct {
	ID               string `json:"id"`
	ShortDescription struct {
		Text string `json:"text"`
	} `json:"shortDescription"`
	Help struct {
		Text string `json:"text"`
	} `json:"help"`
}

type sarifResult struct {
	RuleID  string `json:"ruleId"`
	Level   string `json:"level"`
	Message struct {
		Text string `json:"text"`
	} `json:"message"`
	Locations []sarifLocation `json:"locations,omitempty"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           sarifRegion           `json:"region"`
}

type sarifArtifactLocation struct {
	URI string `json:"uri"`
}

type sarifRegion struct {
	StartLine int `json:"startLine"`
}
