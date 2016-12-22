package synopsis

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"strings"

	"github.com/johnnywidth/synopsis/composer"
	"github.com/johnnywidth/synopsis/repository/compress"
	"github.com/johnnywidth/synopsis/repository/download"
	"github.com/johnnywidth/synopsis/repository/driver"
)

// UpdateAll update all repositories
func (repo *Repository) UpdateAll(config Config) (composer.PackageJSON, error) {
	pm := make(composer.PackageJSON)

	var d driver.Driver
	switch repo.Type {
	case "git", "vcs", "composer":
		d = &driver.Git{URL: repo.URL}
		err := d.Run()
		if err != nil {
			return nil, err
		}

		if len(pm[d.GetName()]) == 0 {
			pm[d.GetName()] = make(map[string]composer.JSONData)
		}

		if config.File.Archive.Directory == "" {
			pm[d.GetName()] = d.GetPackages()
			return pm, nil
		}

		ch := make(chan composer.JSONData)
		HandleCompress(repo.URL, config, d, ch)

		for i := 0; i < len(d.GetPackages()); i++ {
			p := <-ch

			if p.CompressError != nil {
				return nil, p.CompressError
			}

			pm[d.GetName()][p.Version] = p
		}
	case "bitbucket", "git-bitbucket":
		d = &driver.BitBucket{URL: repo.URL}
		return nil, errors.New("Bitbucket driver not supported yet!")
	default:
		return nil, errors.New("Driver was not specified!")
	}

	return pm, nil
}

// HandleCompress handle compress
func HandleCompress(u string, config Config, d driver.Driver, ch chan composer.JSONData) {
	flag := make(chan bool, len(d.GetPackages()))

	for _, jd := range d.GetPackages() {
		go func(jsonData composer.JSONData) {
			flag <- true
			defer func() {
				<-flag
				ch <- jsonData
			}()

			if config.File.Archive.SkipDev && !jsonData.Stability {
				return
			}

			switch config.File.Archive.Format {
			case "zip":
				err := Compress(u, config, d, &jsonData)
				if err != nil {
					jsonData.CompressError = fmt.Errorf("Zip compress failed. %s", err.Error())
				}
			case "tar":
				jsonData.CompressError = fmt.Errorf("Tar format not sopported yet!")
			default:
				jsonData.CompressError = fmt.Errorf("Invalid archive format: %s!", config.File.Archive.Format)
			}
		}(jd)
	}
}

// Compress make archive from repository branch
func Compress(u string, config Config, d driver.Driver, p *composer.JSONData) error {
	// Prepare path and archive path
	gd := &download.Git{
		Name:    p.Name,
		Version: p.Version,
		URL:     u,
		Source:  p.Source,
		DistDir: config.DistDir,
	}
	gd.Prepare()
	// Prepare compression
	cz := compress.Zip{
		SourcePath: gd.SourcePath,
	}
	cz.Prepare()
	// Run download
	if !gd.PathExist && !cz.IsExist {
		err := gd.Run()
		if err != nil {
			return err
		}
	}
	// Run compression
	if gd.PathExist && !cz.IsExist {
		err := cz.Run()
		if err != nil {
			return err
		}
	}

	if cz.IsExist {
		host := config.File.Homepage
		if config.File.Archive.PrefixURL != "" {
			host = config.File.Archive.PrefixURL
		}

		p.Dist = map[string]string{
			"type":      cz.Format,
			"url":       host + "/" + path.Join(config.File.Archive.Directory, strings.Replace(cz.TargetPath, gd.DistDir, "", -1)),
			"reference": d.GetReference(),
			"shasum":    hashFile(cz.TargetPath),
		}
	}

	return nil
}

// FileChunk is chunk
const FileChunk = 8192

func hashFile(f string) string {
	file, _ := os.Open(f)
	defer file.Close()
	info, _ := file.Stat()
	fileSize := info.Size()
	blocks := uint64(math.Ceil(float64(fileSize) / float64(FileChunk)))
	hash := sha1.New()
	for i := uint64(0); i < blocks; i++ {
		blockSize := int(math.Min(FileChunk, float64(fileSize-int64(i*FileChunk))))
		buf := make([]byte, blockSize)
		file.Read(buf)
		io.WriteString(hash, string(buf))
	}
	return hex.EncodeToString(hash.Sum(nil))
}
