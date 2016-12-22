package composer

// JSONSchema composer json schema
// https://getcomposer.org/doc/04-schema.md#json-schema
type JSONSchema struct {
	Name               string                 `json:"name,omitempty"`
	Description        string                 `json:"description,omitempty"`
	Version            string                 `json:"version,omitempty"`
	Type               string                 `json:"type,omitempty"`
	Keywords           []string               `json:"keywords,omitempty"`
	Homepage           string                 `json:"homepage,omitempty"`
	Time               string                 `json:"time,omitempty"`
	License            interface{}            `json:"license,omitempty"`
	Authors            []map[string]string    `json:"authors,omitempty"`
	Support            map[string]string      `json:"support,omitempty"`
	Require            map[string]string      `json:"require,omitempty"`
	RequireDev         map[string]string      `json:"require-dev,omitempty"`
	Conflict           interface{}            `json:"conflict,omitempty"`
	Replace            interface{}            `json:"replace,omitempty"`
	Provide            interface{}            `json:"provide,omitempty"`
	Suggest            interface{}            `json:"suggest,omitempty"`
	Autoload           map[string]interface{} `json:"autoload,omitempty"`
	AutoloadDev        map[string]interface{} `json:"autoload-dev,omitempty"`
	IncludePath        interface{}            `json:"include-path,omitempty"`
	TargetDir          interface{}            `json:"target-dir,omitempty"`
	MinimumStability   interface{}            `json:"minimum-stability,omitempty"`
	PreferStable       interface{}            `json:"prefer-stable,omitempty"`
	Repositories       interface{}            `json:"repositories,omitempty"`
	Config             interface{}            `json:"config,omitempty"`
	Scripts            interface{}            `json:"scripts,omitempty"`
	Extra              map[string]interface{} `json:"extra,omitempty"`
	Bin                interface{}            `json:"bin,omitempty"`
	Archive            interface{}            `json:"archive,omitempty"`
	NonFeatureBranches interface{}            `json:"non-feature-branches,omitempty"`
}
