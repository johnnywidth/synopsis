package composer

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path"
)

// Packages is a packages field from packages.json
type Packages struct {
	Package PackageJSON `json:"packages"`
}

// PackageJSON package json
type PackageJSON map[string]map[string]JSONData

// JSONData data from repository composer.json and data for packages.josn
type JSONData struct {
	JSONSchema

	VersionNormalized string            `json:"version_normalized,omitempty"`
	Source            map[string]string `json:"source,omitempty"`
	Dist              map[string]string `json:"dist,omitempty"`
	InstallSource     string            `json:"installation-source,omitempty"`

	Stability     bool  `json:"-"`
	CompressError error `json:"-"`
}

// ComposerCacheVcs path to composer cache
const ComposerCacheVcs = "/.composer/cache/vcs/"

// InitPackages init packages
func InitPackages() error {
	ccv := path.Join(os.Getenv("HOME"), ComposerCacheVcs)

	err := os.MkdirAll(ccv, 0777)

	return err
}

// WriteToFile crete json from packages and write to file
func (packages *Packages) WriteToFile(output string) error {
	j, err := json.MarshalIndent(packages, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(path.Join(output, "packages.json"), j, 0755)
}
