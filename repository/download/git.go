package download

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
)

// Git git repository data
type Git struct {
	Name       string
	Version    string
	URL        string
	Source     map[string]string
	DistDir    string
	SourcePath string
	PathExist  bool
}

var pRegExp = regexp.MustCompile("[^a-z0-9-_]")

// Prepare make prepare before download
func (git *Git) Prepare() {
	dir := pRegExp.ReplaceAllString(git.Name, "-")

	h := sha1.New()
	h.Write([]byte(git.Source["reference"]))

	ref := hex.EncodeToString(h.Sum(nil))
	ref = ref[:6]

	temp := dir + "-" + git.Version + "-" + ref
	git.SourcePath = path.Join(git.DistDir, strings.Replace(temp, "/", "-", -1))
	git.PathExist = false

	_, err := os.Stat(git.SourcePath)
	if err == nil {
		git.PathExist = true
	}
}

// Run run git downloader
func (git *Git) Run() error {
	cmd := exec.Command("git", "clone", "--no-checkout", git.URL, git.SourcePath)

	_, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git clone --no-checkout %s %s. %s", git.URL, git.SourcePath, err)
	}

	cmd = exec.Command("git", "remote", "add", "composer", git.URL)
	cmd.Dir = git.SourcePath

	_, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git remote add composer %s. %s", git.URL, err)
	}

	cmd = exec.Command("git", "fecth", "composer")
	cmd.Dir = git.SourcePath
	cmd.CombinedOutput()

	cmd = exec.Command("git", "checkout", "-b", git.Version, "composer/"+git.Version)
	cmd.Dir = git.SourcePath
	cmd.CombinedOutput()

	cmd = exec.Command("git", "reset", "--hard", git.Source["reference"])
	cmd.Dir = git.SourcePath
	cmd.CombinedOutput()

	_, err = os.Stat(git.SourcePath)
	if err != nil {
		return err
	}

	git.PathExist = true
	return nil
}
