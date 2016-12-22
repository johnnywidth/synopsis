package driver

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/johnnywidth/synopsis/composer"
)

// Git git driver data
type Git struct {
	Name              string
	RepoDir           string
	Version           string
	VersionNormalized string
	Reference         string
	Source            map[string]string
	Dist              map[string]string
	URL               string
	Packages          map[string]composer.JSONData
}

// Run run git driver
func (git *Git) Run() error {
	err := git.prepareRepoDir()
	if err != nil {
		return err
	}

	err = git.PrepareMainBranch()
	if err != nil {
		return err
	}

	git.Packages = make(map[string]composer.JSONData)

	err = git.PrepareTags()
	if err != nil {
		return err
	}

	err = git.PrepareBranches()
	if err != nil {
		return err
	}

	return nil
}

// GetName get repository name
func (git *Git) GetName() string {
	return git.Name
}

// GetSource get repository source
func (git *Git) GetSource() map[string]string {
	return git.Source
}

// GetReference get repository reference
func (git *Git) GetReference() string {
	return git.Reference
}

// GetPackages get composer
func (git *Git) GetPackages() map[string]composer.JSONData {
	return git.Packages
}

var pmbRegExp = regexp.MustCompile("\\* +(\\S+)")

// PrepareMainBranch prepare main branch
// Set to driver `Reference`, `Name`
func (git *Git) PrepareMainBranch() error {
	cmd := exec.Command("git", "branch", "--no-color")
	cmd.Dir = git.RepoDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git branch --no-color %s", err)
	}

	response := pmbRegExp.FindStringSubmatch(string(out))
	if response == nil || len(response) < 2 {
		return fmt.Errorf("Prepare main branch failed. %s", git.URL)
	}

	git.Reference = response[1]

	composer, err := git.getComposerInformation()
	if err != nil {
		return err
	}

	git.Name = composer.Name

	return nil
}

var ptRegExp = regexp.MustCompile("^([a-f0-9]{40}) refs/tags/(\\S+)")

// PrepareTags prepare repository tags
// Set to driver `Version`, `VersionNormalized`, `Reference`
// Set composer json data by version
func (git *Git) PrepareTags() error {
	cmd := exec.Command("git", "show-ref", "--tags")
	cmd.Dir = git.RepoDir

	// Don't need check error, because tags may not exist
	out, _ := cmd.CombinedOutput()

	for _, tag := range strings.SplitAfter(string(out), "\n") {
		response := ptRegExp.FindStringSubmatch(tag)
		if response == nil {
			continue
		}

		if response != nil && len(response) < 3 {
			return fmt.Errorf("Prepare tags failed. %s", git.URL)
		}

		git.Version = composer.PrepareTagVersion(response[2])
		newVersion := composer.VersionNormalizedTag(response[2])
		git.VersionNormalized = composer.PrepareTagVersionNormalized(newVersion)
		git.Reference = response[1]

		jsonData, err := git.getComposerInformation()
		if err != nil {
			return fmt.Errorf("Prepare tags failed. %s. Get composer information. %s", git.URL, err.Error())
		}

		stability := composer.GetStability(git.Version)
		if stability != composer.Dev {
			jsonData.Stability = true
		}

		git.Packages[git.Version] = jsonData
	}

	return nil
}

var pbRegExp = regexp.MustCompile("(?:\\* )? *(\\S+) *([a-f0-9]+)(?: .*)?")

// PrepareBranches prepare repository branches
// Set to driver `Version`, `VersionNormalized`, `Reference`
// Set composer json data by version
func (git *Git) PrepareBranches() error {
	cmd := exec.Command("git", "branch", "--no-color", "--no-abbrev", "-v")
	cmd.Dir = git.RepoDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Prepare branches failed. %s. git branches --no-color --no-abbrev -v. %s", git.URL, err.Error())
	}

	for _, branch := range strings.SplitAfter(string(out), "\n") {
		response := pbRegExp.FindStringSubmatch(branch)
		if response == nil {
			continue
		}

		if response != nil && len(response) < 3 {
			return fmt.Errorf("Prepare branches failed. %s", git.URL)
		}

		git.Version = composer.PrepareBranchVersion(response[1])
		git.VersionNormalized = composer.VersionNormalizedBranch(response[1])
		git.Reference = response[2]

		jsonData, err := git.getComposerInformation()
		if err != nil {
			return fmt.Errorf("Prepare branches failed. %s. Get composer information. %s", git.URL, err.Error())
		}

		stability := composer.GetStability(git.Version)
		if stability != "dev" {
			jsonData.Stability = true
		}

		git.Packages[git.Version] = jsonData
	}

	return nil
}

var prdRegExp = regexp.MustCompile("[^a-z0-9.]")

// prepareRepoDir prepare `RepoDir` for driver
func (git *Git) prepareRepoDir() error {
	dir := prdRegExp.ReplaceAllString(git.URL, "-")

	git.RepoDir = os.Getenv("HOME") + composer.ComposerCacheVcs + dir
	_, err := os.Stat(git.RepoDir)
	if err != nil {
		cmd := exec.Command("git", "clone", "--mirror", git.URL, git.RepoDir)
		cmd.Dir = os.Getenv("HOME") + composer.ComposerCacheVcs

		_, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git clone --mirror %s %s. %s", git.URL, git.RepoDir, err)
		}
	} else {
		cmd := exec.Command("git", "remote", "set-url", "origin", git.URL)
		cmd.Dir = git.RepoDir

		_, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git remote set-url origin %s. %s", git.URL, err)
		}
		cmd = exec.Command("git", "remote", "update", "--prune", "origin")
		cmd.Dir = git.RepoDir

		_, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git remote update --prune origin. %s. %s", git.URL, err)
		}
	}

	return nil
}

func (git *Git) getComposerInformation() (composer.JSONData, error) {
	cmd := exec.Command("git", "show", git.Reference+":composer.json")
	cmd.Dir = git.RepoDir

	out, err := cmd.CombinedOutput()
	if err != nil {
		return composer.JSONData{}, fmt.Errorf("git show %s. %s", git.Reference+":composer.json", err)
	}

	jsonData := composer.JSONData{}

	err = json.Unmarshal(out, &jsonData)
	if err != nil {
		return composer.JSONData{}, err
	}

	if jsonData.Time == "" {
		cmd = exec.Command("git", "log", "-1", "--format=%at", git.Reference)
		cmd.Dir = git.RepoDir

		out, err := cmd.CombinedOutput()
		if err != nil {
			return composer.JSONData{}, fmt.Errorf("git log -1 --format. %s", err)
		}

		t, _ := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
		jsonData.Time = time.Unix(t, 0).String()
	}

	git.Source = map[string]string{
		"type":      "git",
		"url":       git.URL,
		"reference": git.Reference,
	}

	jsonData.Source = git.Source
	jsonData.Version = git.Version
	jsonData.VersionNormalized = git.VersionNormalized

	return jsonData, nil
}
