package diagnostic

type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

type Item struct {
	RuleID   string   `json:"rule_id"`
	Severity Severity `json:"severity"`
	File     string   `json:"file,omitempty"`
	Line     int      `json:"line,omitempty"`
	Message  string   `json:"message"`
	Why      string   `json:"why,omitempty"`
	FixHint  string   `json:"fix_hint,omitempty"`
}

func (i Item) IsError() bool {
	return i.Severity == SeverityError
}
