package compress

import (
	"archive/zip"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Zip compression data
type Zip struct {
	Format     string
	TargetPath string
	SourcePath string
	IsExist    bool
}

// Prepare make prepare before compression
func (a *Zip) Prepare() {
	a.Format = "zip"
	a.TargetPath = a.SourcePath + "." + a.Format

	_, err := os.Stat(a.TargetPath)
	if err == nil {
		a.IsExist = true
	}
}

// Run start compression
func (a *Zip) Run() error {
	zipFile, err := os.Create(a.TargetPath)
	if err != nil {
		return err
	}

	zipWriter := zip.NewWriter(zipFile)

	defer func() {
		zipWriter.Close()
		os.RemoveAll(a.SourcePath)
	}()

	err = filepath.Walk(a.SourcePath, func(p string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			buffer, err := ioutil.ReadFile(p)
			if err != nil {
				return err
			}

			file := strings.Trim(strings.Replace(p, a.SourcePath, "", -1), "/")

			f, err := zipWriter.Create(file)
			if err != nil {
				return err
			}
			if _, err = f.Write([]byte(buffer)); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	a.IsExist = true
	return nil
}
