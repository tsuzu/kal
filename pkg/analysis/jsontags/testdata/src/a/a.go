package a

type JSONTagTestStruct struct {
	NoJSONTag            string // want "field NoJSONTag is missing json tag"
	EmptyJSONTag         string `json:""`                        // want "field EmptyJSONTag has empty json tag"
	NonCamelCaseJSONTag  string `json:"non_camel_case_json_tag"` // want "field NonCamelCaseJSONTag has non-camel case json tag: non_camel_case_json_tag"
	WithHyphensJSONTag   string `json:"with-hyphens-json-tag"`   // want "field WithHyphensJSONTag has non-camel case json tag: with-hyphens-json-tag"
	NonTerminatedJSONTag string `json:"nonTerminatedJSONTag`     // want "field NonTerminatedJSONTag is missing json tag"
	XMLTag               string `xml:"xmlTag"`                   // want "field XMLTag is missing json tag"
	InlineJSONTag        string `json:",inline"`
	ValidJSONTag         string `json:"validJsonTag"`
	ValidOptionalJSONTag string `json:"validOptionalJsonTag,omitempty"`
	JSONTagWithID        string `json:"jsonTagWithID"`
	JSONTagWithTTL       string `json:"jsonTagWithTTL"`
	JSONTagWithGiB       string `json:"jsonTagWithGiB"`
}
