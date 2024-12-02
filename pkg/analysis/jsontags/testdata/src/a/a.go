package a

type JSONTagTestStruct struct {
	NoJSONTag            string // want "field NoJSONTag is missing json tag"
	EmptyJSONTag         string `json:""`                        // want "field EmptyJSONTag has empty json tag"
	NonCamelCaseJSONTag  string `json:"non_camel_case_json_tag"` // want "field NonCamelCaseJSONTag json tag does not match pattern \"\\^\\[a-z\\]\\[a-z0-9\\]\\*\\(\\?:\\[A-Z\\]\\[a-z0-9\\]\\*\\)\\*\\$\": non_camel_case_json_tag"
	WithHyphensJSONTag   string `json:"with-hyphens-json-tag"`   // want "field WithHyphensJSONTag json tag does not match pattern \"\\^\\[a-z\\]\\[a-z0-9\\]\\*\\(\\?:\\[A-Z\\]\\[a-z0-9\\]\\*\\)\\*\\$\": with-hyphens-json-tag"
	PascalCaseJSONTag    string `json:"PascalCaseJSONTag"`       // want "field PascalCaseJSONTag json tag does not match pattern \"\\^\\[a-z\\]\\[a-z0-9\\]\\*\\(\\?:\\[A-Z\\]\\[a-z0-9\\]\\*\\)\\*\\$\": PascalCaseJSONTag"
	NonTerminatedJSONTag string `json:"nonTerminatedJSONTag`     // want "field NonTerminatedJSONTag is missing json tag"
	XMLTag               string `xml:"xmlTag"`                   // want "field XMLTag is missing json tag"
	InlineJSONTag        string `json:",inline"`
	ValidJSONTag         string `json:"validJsonTag"`
	ValidOptionalJSONTag string `json:"validOptionalJsonTag,omitempty"`
	JSONTagWithID        string `json:"jsonTagWithID"`
	JSONTagWithTTL       string `json:"jsonTagWithTTL"`
	JSONTagWithGiB       string `json:"jsonTagWithGiB"`
}
