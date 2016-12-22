package composer_test

import (
	"testing"

	"github.com/johnnywidth/synopsis/composer"
)

type TagTestCase struct {
	Tag               string
	Version           string
	NewNormalized     string
	VersionNormalized string
	Stability         string
}

type BranchTestCase struct {
	Branch            string
	Version           string
	VersionNormalized string
	Stability         string
}

var ttc []TagTestCase
var btc []BranchTestCase

func init() {
	ttc = []TagTestCase{{
		Tag:               "0",
		Version:           "0",
		NewNormalized:     "0.0.0.0",
		VersionNormalized: "0.0.0.0",
		Stability:         "stable",
	}, {
		Tag:               "1",
		Version:           "1",
		NewNormalized:     "1.0.0.0",
		VersionNormalized: "1.0.0.0",
		Stability:         "stable",
	}, {
		Tag:               "1.0.0",
		Version:           "1.0.0",
		NewNormalized:     "1.0.0.0",
		VersionNormalized: "1.0.0.0",
		Stability:         "stable",
	}, {
		Tag:               "v0.2.0",
		Version:           "v0.2.0",
		NewNormalized:     "0.2.0.0",
		VersionNormalized: "0.2.0.0",
		Stability:         "stable",
	}, {
		Tag:               "v1.0.0-RC",
		Version:           "v1.0.0-RC",
		NewNormalized:     "1.0.0.0-RC",
		VersionNormalized: "1.0.0.0-RC",
		Stability:         "RC",
	}, {
		Tag:               "v1.0.0-RC1",
		Version:           "v1.0.0-RC1",
		NewNormalized:     "1.0.0.0-RC1",
		VersionNormalized: "1.0.0.0-RC1",
		Stability:         "RC",
	}, {
		Tag:               "v1.0.0-b",
		Version:           "v1.0.0-b",
		NewNormalized:     "1.0.0.0-beta",
		VersionNormalized: "1.0.0.0-beta",
		Stability:         "beta",
	}, {
		Tag:               "v1.0.0-a",
		Version:           "v1.0.0-a",
		NewNormalized:     "1.0.0.0-alpha",
		VersionNormalized: "1.0.0.0-alpha",
		Stability:         "alpha",
	}, {
		Tag:               "v1.0.0-p",
		Version:           "v1.0.0-p",
		NewNormalized:     "1.0.0.0-patch",
		VersionNormalized: "1.0.0.0-patch",
		Stability:         "stable",
	}, {
		Tag:               "v1.2.0-alpha",
		Version:           "v1.2.0-alpha",
		NewNormalized:     "1.2.0.0-alpha",
		VersionNormalized: "1.2.0.0-alpha",
		Stability:         "alpha",
	}, {
		Tag:               "v1.2.0-alpha2",
		Version:           "v1.2.0-alpha2",
		NewNormalized:     "1.2.0.0-alpha2",
		VersionNormalized: "1.2.0.0-alpha2",
		Stability:         "alpha",
	},
		//{
		//	Tag:               "",
		//	Version:           "1.0.0RC1dev",
		//	NewNormalized:     "1.0.0.0-RC1-dev",
		//	VersionNormalized: "1.0.0.0-RC1",
		//	Stability:         "dev",
		//},
		//{
		//	Tag:               "1.0.0-rC15-dev",
		//	Version:           "1.0.0-rC15-dev",
		//	NewNormalized:     "1.0.0.0-RC15-dev",
		//	VersionNormalized: "1.0.0.0-RC15",
		//	Stability:         "dev",
		//},
		{
			Tag:               "1.0.0-rc1",
			Version:           "1.0.0-rc1",
			NewNormalized:     "1.0.0.0-RC1",
			VersionNormalized: "1.0.0.0-RC1",
			Stability:         "RC",
		}, {
			Tag:               "10.4.13beta.2",
			Version:           "10.4.13beta.2",
			NewNormalized:     "10.4.13.0-beta2",
			VersionNormalized: "10.4.13.0-beta2",
			Stability:         "beta",
		},
		//{
		//	Tag:               "1.0.0-beta.5+foo",
		//	Version:           "1.0.0-beta.5+foo",
		//	NewNormalized:     "1.0.0.0-beta5",
		//	VersionNormalized: "1.0.0.0-beta5",
		//	Stability:         "beta",
		//},
		{
			Tag:               "1.0.0+foo",
			Version:           "1.0.0+foo",
			NewNormalized:     "1.0.0.0",
			VersionNormalized: "1.0.0.0",
			Stability:         "stable",
		}, {
			Tag:               "1.0.0+foo as 2.0",
			Version:           "1.0.0+foo as 2.0",
			NewNormalized:     "1.0.0.0",
			VersionNormalized: "1.0.0.0",
			Stability:         "stable",
		}}

	btc = []BranchTestCase{{
		Branch:            "master",
		Version:           "dev-master",
		VersionNormalized: "9999999-dev",
		Stability:         "dev",
	}, {
		Branch:            "dev",
		Version:           "dev-dev",
		VersionNormalized: "dev-dev",
		Stability:         "dev",
	}, {
		Branch:            "develop-2.3.0",
		Version:           "dev-develop-2.3.0",
		VersionNormalized: "dev-develop-2.3.0",
		Stability:         "dev",
	}, {
		Branch:            "pr/171",
		Version:           "dev-pr/171",
		VersionNormalized: "dev-pr/171",
		Stability:         "dev",
	}, {
		Branch:            "revert-354-master",
		Version:           "dev-revert-354-master",
		VersionNormalized: "dev-revert-354-master",
		Stability:         "dev",
	}, {
		Branch:            "v2.0",
		Version:           "2.0.x-dev",
		VersionNormalized: "2.0.9999999.9999999-dev",
		Stability:         "dev",
	}, {
		Branch:            "issue/126",
		Version:           "dev-issue/126",
		VersionNormalized: "dev-issue/126",
		Stability:         "dev",
	}, {
		Branch:            "dev-1.7",
		Version:           "dev-dev-1.7",
		VersionNormalized: "dev-dev-1.7",
		Stability:         "dev",
	}}
}

func TestTagVersion(t *testing.T) {
	for _, tc := range ttc {
		v := composer.PrepareTagVersion(tc.Tag)
		if v != tc.Version {
			t.Fatalf("PrepareTagVersion: %s != %s", v, tc.Version)
		}

		nv := composer.VersionNormalizedTag(tc.Tag)
		if nv != tc.NewNormalized {
			t.Fatalf("VersionNormalizedTag: %s != %s", nv, tc.NewNormalized)
		}

		vn := composer.PrepareTagVersionNormalized(tc.NewNormalized)
		if vn != tc.VersionNormalized {
			t.Fatalf("PrepareTagVersionNormalized: %s != %s", vn, tc.VersionNormalized)
		}

		s := composer.GetStability(tc.Version)
		if s != tc.Stability {
			t.Fatalf("GetStability: %s != %s", s, tc.Stability)
		}
	}
}

func TestBranchVersion(t *testing.T) {
	for _, tc := range btc {
		v := composer.PrepareBranchVersion(tc.Branch)
		if v != tc.Version {
			t.Fatalf("PrepareBranchVersion: %s != %s", v, tc.Version)
		}

		vn := composer.VersionNormalizedBranch(tc.Branch)
		if vn != tc.VersionNormalized {
			t.Fatalf("VersionNormalizedBranch: %s != %s", vn, tc.VersionNormalized)
		}

		s := composer.GetStability(tc.Version)
		if s != tc.Stability {
			t.Fatalf("GetStability: %s != %s", s, tc.Stability)
		}
	}
}

func BenchmarkPrepareTagVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tCase := range ttc {
			composer.PrepareTagVersion(tCase.Tag)
		}
	}
}

func BenchmarkVersionNormalizedTag(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tCase := range ttc {
			composer.VersionNormalizedTag(tCase.Tag)
		}
	}
}

func BenchmarkPrepareTagVersionNormalized(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tCase := range ttc {
			composer.PrepareTagVersionNormalized(tCase.NewNormalized)
		}
	}
}

func BenchmarkPrepareBranchVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tCase := range btc {
			composer.PrepareBranchVersion(tCase.Branch)
		}
	}
}

func BenchmarkVersionNormalizedBranch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tCase := range btc {
			composer.VersionNormalizedBranch(tCase.Branch)
		}
	}
}

func BenchmarkGetStability(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, tCase := range ttc {
			composer.GetStability(tCase.Version)
		}

		for _, tCase := range btc {
			composer.GetStability(tCase.Version)
		}
	}
}
