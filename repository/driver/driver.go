package driver

import "github.com/johnnywidth/synopsis/composer"

// Driver is interface for VCS drivers
type Driver interface {
	Run() error
	GetName() string
	GetSource() map[string]string
	GetReference() string
	GetPackages() map[string]composer.JSONData

	// PrepareMainBranch prepare main branch
	// Set to driver `Reference`, `Name`
	PrepareMainBranch() error

	// PrepareTags prepare repository tags
	// Set to driver `Version`, `VersionNormalized`, `Reference`
	// Set composer json data by version
	PrepareTags() error

	// PrepareBranches prepare repository branches
	// Set to driver `Version`, `VersionNormalized`, `Reference`
	// Set composer json data by version
	PrepareBranches() error
}
