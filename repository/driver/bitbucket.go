package driver

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"

	"github.com/johnnywidth/synopsis/composer"
)

// BitBucket bitbucket driver data
type BitBucket struct {
	Owner      string
	Repository string
	Version    string
	Reference  string
	Source     map[string]string
	Dist       map[string]string
	URL        string
	Packages   map[string]composer.JSONData
}

// Run run bitbucket driver
func (bucket *BitBucket) Run() error {
	bucket.prepareRepository()
	if err := bucket.PrepareMainBranch(); err != nil {
		return err
	}
	return nil
}

// GetName get repository name
func (bucket *BitBucket) GetName() string {
	return bucket.Owner
}

// GetSource get repository source
func (bucket *BitBucket) GetSource() map[string]string {
	return bucket.Source
}

// GetReference get repository reference
func (bucket *BitBucket) GetReference() string {
	return bucket.Reference
}

// GetPackages get composer
func (bucket *BitBucket) GetPackages() map[string]composer.JSONData {
	return bucket.Packages
}

// PrepareMainBranch prepare main branch
// Set to driver `Reference`, `Name`
func (bucket *BitBucket) PrepareMainBranch() error {
	url := "https://api.bitbucket.org/1.0/repositories/" + bucket.Owner + "/" + bucket.Repository
	var r interface{}
	_, err := httpGet(url, r)
	if err != nil {
		return err
	}
	return nil
}

func (bucket *BitBucket) getComposerInformation(commit string) composer.JSONData {
	return composer.JSONData{}
}

// PrepareTags prepare repository tags
// Set to driver `Version`, `VersionNormalized`, `Reference`
// Set composer json data by version
func (bucket *BitBucket) PrepareTags() error {
	return nil
}

// PrepareBranches prepare repository branches
// Set to driver `Version`, `VersionNormalized`, `Reference`
// Set composer json data by version
func (bucket *BitBucket) PrepareBranches() error {
	return nil
}

var prRegExp = regexp.MustCompile("^https?://bitbucket\\.org/([^/]+)/(.+?)\\.git$")

func (bucket *BitBucket) prepareRepository() {
	response := prRegExp.FindStringSubmatch(bucket.URL)
	bucket.Owner, bucket.Repository = response[1], response[2]
}

func (bucket *BitBucket) prepareSource(commit string) {

}

func httpGet(url string, result interface{}) (interface{}, error) {
	response, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	json.Unmarshal(body, &result)
	return result, nil
}
