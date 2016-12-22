package synopsis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// Config is app config
type Config struct {
	File         File
	FileName     string
	ThreadNumber int
	OutputDir    string
	DistDir      string
}

// File is config file
type File struct {
	Name         string       `json:"name"`
	Homepage     string       `json:"homepage"`
	Archive      Archive      `json:"archive"`
	Repositories []Repository `json:"repositories"`
}

// Archive is archive field in config file
type Archive struct {
	Directory string `json:"directory"`
	Format    string `json:"format"`
	SkipDev   bool   `json:"skip-dev"`
	PrefixURL string `json:"prefix-url"`
}

// Repository is repository field in config file
type Repository struct {
	Type string `json:"type"`
	URL  string `json:"url"`
}

// PrepareConfig read config file and create config structure
func (config *Config) PrepareConfig(cf string, od string, tn string) error {
	config.FileName = cf

	file, err := ioutil.ReadFile(config.FileName)
	if err != nil {
		return fmt.Errorf("Wrong app config. `FILE`. %s", err.Error())
	}

	err = json.Unmarshal(file, &config.File)
	if err != nil {
		return fmt.Errorf("Unmarshall failed. %s", err.Error())
	}

	config.ThreadNumber, err = strconv.Atoi(tn)
	if err != nil {
		return fmt.Errorf("Wrong app config. `THREAD`. %s", err.Error())
	}

	config.OutputDir = od
	if config.OutputDir == "" {
		return fmt.Errorf("Wrong app config. `OUTPUT` is empty!")
	}

	if config.File.Archive.Directory == "" {
		return fmt.Errorf("Config param `archive.directory` is empty!")
	}

	config.DistDir = path.Join(config.OutputDir, config.File.Archive.Directory)

	if config.File.Archive.Format == "" {
		config.File.Archive.Format = "zip"
	}

	return nil
}

// MakeOutputDir create dir for archive repository branch and packages.json file
func (config *Config) MakeOutputDir() {
	_, err := os.Stat(config.OutputDir)
	if err != nil {
		os.MkdirAll(config.OutputDir, 0777)
	}
}
