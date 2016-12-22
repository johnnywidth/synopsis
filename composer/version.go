package composer

import (
	"regexp"
	"strings"
)

var modifier = `[._-]?(?:(stable|beta|b|RC|alpha|a|patch|pl|p)(?:[.-]?(\d+))?)?([.-]?dev)?`

var vntRegExp1 = regexp.MustCompile(`(?i)^v?(\d{1,3})(\.\d+)?(\.\d+)?(\.\d+)?` + modifier)
var vntRegExp2 = regexp.MustCompile("(?i)[.-]?dev$")

// DevPrefix dev
const DevPrefix = "dev-"

// VersionNormalizedTag normalize tag
func VersionNormalizedTag(v string) string {
	v = strings.TrimSpace(v)
	if len(v) >= 4 && DevPrefix == strings.ToLower(v[:4]) {
		return DevPrefix + v[4:]
	}
	// match classical number version
	resp1 := vntRegExp1.FindStringSubmatch(v)
	if resp1 != nil && len(resp1) >= 5 {
		var vers string
		index := 5
		for i := 1; i < index; i++ {
			if resp1[i] != "" {
				vers += resp1[i]
			} else {
				vers += ".0"
			}
		}

		if resp1[index] == "" || resp1[index] == "stable" {
			return vers
		}

		vers += "-" + expandStability(resp1[index])
		if len(resp1) >= index+1 && resp1[index+1] != "" {
			vers += resp1[index+1]
		}
		if len(resp1) >= index+2 && resp1[index+2] != "" {
			vers += "-dev"
		}

		return vers
	}
	// match dev branches
	resp2 := vntRegExp2.FindStringSubmatch(v)
	if resp2 != nil && len(resp2) >= 2 {
		return VersionNormalizedBranch(resp2[1])
	}
	// return default version
	return v
}

var vnbRegExp1 = regexp.MustCompile(`(?i)^(?:dev-)?(?:master|trunk|default)$`)
var vnbRegExp2 = regexp.MustCompile(`(?i)^v?(\d+)(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?(\.(?:\d+|[xX*]))?$`)

// VersionNormalizedBranch normalize branch
func VersionNormalizedBranch(v string) string {
	v = strings.TrimSpace(v)
	// match master-like branch
	resp1 := vnbRegExp1.FindStringSubmatch(v)
	if resp1 != nil {
		return "9999999-dev"
	}
	// match dev branches
	resp2 := vnbRegExp2.FindStringSubmatch(v)
	// return default version
	if resp2 == nil || len(resp2) < 5 {
		return DevPrefix + v
	}

	var nv string
	for i := 1; i < 5; i++ {
		if resp2[i] != "" {
			resp2[i] = strings.Replace(resp2[i], "*", "x", -1)
			resp2[i] = strings.Replace(resp2[i], "X", "x", -1)
			nv += resp2[i]
		} else {
			nv += ".x"
		}
	}

	return strings.Replace(nv, "x", "9999999", -1) + "-dev"
}

var ptvRegExp = regexp.MustCompile("(?i)(.*?)[.-]?dev$")

// PrepareTagVersion prepare tag
func PrepareTagVersion(v string) string {
	return ptvRegExp.ReplaceAllString(v, "")
}

var ptvnRegExp = regexp.MustCompile("(?i)(^dev-|[.-]?dev$)")

// PrepareTagVersionNormalized prepare tag normalized
func PrepareTagVersionNormalized(v string) string {
	return ptvnRegExp.ReplaceAllString(v, "")
}

var pbvRegExp = regexp.MustCompile(`(\.9{7})+`)

// PrepareBranchVersion prepare branch
func PrepareBranchVersion(v string) string {
	newV := VersionNormalizedBranch(v)
	if DevPrefix == newV[:4] || "9999999-dev" == newV {
		return DevPrefix + v
	}

	return pbvRegExp.ReplaceAllString(newV, ".x")
}

var isRegExp1 = regexp.MustCompile(`(?i)#.+$`)
var isRegExp2 = regexp.MustCompile(`(?i)` + modifier + `$`)

// Realeases
const (
	Stable = "stable"
	Dev    = "dev"
	Beta   = "beta"
	Alpha  = "alpha"
	RC     = "RC"
)

// GetStability return stability
func GetStability(v string) string {
	v = isRegExp1.ReplaceAllString(v, "")

	if len(v) >= 4 && (DevPrefix == v[:4] || "-dev" == v[len(v)-4:]) {
		return Dev
	}

	resp := isRegExp2.FindStringSubmatch(strings.ToLower(v))
	if resp != nil {
		if len(resp) >= 4 && resp[3] != "" {
			return Dev
		}

		if len(resp) >= 2 {
			if resp[1] == Beta || resp[1] == "b" {
				return Beta
			}

			if resp[1] == Alpha || resp[1] == "a" {
				return Alpha
			}

			if resp[1] == strings.ToLower(RC) {
				return RC
			}
		}
	}

	return Stable
}

func expandStability(stability string) string {
	stability = strings.ToLower(stability)
	switch stability {
	case "a":
		return Alpha
	case "b":
		return Beta
	case "p", "pl":
		return "patch"
	case "rc":
		return RC
	}

	return stability
}
