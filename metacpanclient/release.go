package metacpanclient

import (
	"encoding/json"
	"strings"

	// local
	cm "github.com/cmburn/perlutils/cpanmeta"
	pui "github.com/cmburn/perlutils/internal"
	"github.com/cmburn/perlutils/version"
)

type Release struct {
	Dependency []struct {
		Phase        Phase        `json:"phase"`
		Relationship Relationship `json:"relationship"`
		Module       string       `json:"module"`
		Version      version.JSON `json:"version"`
	} `json:"dependency"`
	Stat struct {
		MTime int `json:"mtime"`
		Mode  int `json:"mode"`
		Size  int `json:"size"`
	} `json:"stat"`
	Tests struct {
		Unknown int `json:"unknown"`
		Pass    int `json:"pass"`
		Fail    int `json:"fail"`
		NA      int `json:"na"`
	} `json:"tests"`
	Abstract        string        `json:"abstract"`
	Archive         string        `json:"archive"`
	Author          string        `json:"author"`
	ChecksumMD5     string        `json:"checksum_md5"`
	ChecksumSHA256  string        `json:"checksum_sha256"`
	Date            string        `json:"date"`
	Distribution    string        `json:"Distribution"`
	DownloadURL     string        `json:"download_url"`
	License         []string      `json:"license"`
	MainModule      string        `json:"main_module"`
	Maturity        Maturity      `json:"maturity"`
	Metadata        cm.Spec       `json:"metadata"`
	Name            string        `json:"name"`
	Provides        []string      `json:"provides"`
	Resources       cm.Resources  `json:"resources"`
	Status          ReleaseStatus `json:"status"`
	Version         version.JSON  `json:"version"`
	VersionNumified float64       `json:"version_numified"`
	Authorized      bool          `json:"authorized"`
	Deprecated      bool          `json:"deprecated"`
	First           bool          `json:"first"`
	hasUA
	changes string
}

func (r *Release) Changes() (string, error) {
	if r.changes != "" {
		return r.changes, nil
	}
	if r.mc == nil {
		return "", ErrNilClient
	}
	sb := strings.Builder{}
	sb.WriteString("changes/")
	sb.WriteString(r.Author)
	sb.WriteRune('/')
	sb.WriteString(r.Name)
	req := NewRequest[*wrapper[struct {
		Content string `json:"content"`
	}]]("", "", r.mc.Debug(), r.mc)
	c, err := req.Fetch(sb.String(), nil)
	if err != nil {
		return "", err
	}
	r.changes = c.Result.Content
	return r.changes, nil
}

func (r *Release) MetaCPANURL() string {
	const baseURL = "https://metacpan.org/release/"
	sb := strings.Builder{}
	sb.WriteString(baseURL)
	sb.WriteString(r.Author)
	sb.WriteRune('/')
	sb.WriteString(r.Name)
	return sb.String()
}

func (r *Release) UnmarshalJSON(data []byte) error {
	var v struct {
		Dependency []struct {
			Phase        Phase        `json:"phase"`
			Relationship Relationship `json:"relationship"`
			Module       string       `json:"module"`
			Version      version.JSON `json:"version"`
		} `json:"dependency"`
		Stat struct {
			MTime pui.CoercedInt `json:"mtime"`
			Mode  pui.CoercedInt `json:"mode"`
			Size  pui.CoercedInt `json:"size"`
		} `json:"stat"`
		Tests struct {
			Unknown pui.CoercedInt `json:"unknown"`
			Pass    pui.CoercedInt `json:"pass"`
			Fail    pui.CoercedInt `json:"fail"`
			NA      pui.CoercedInt `json:"na"`
		} `json:"tests"`
		Abstract        string        `json:"abstract"`
		Archive         string        `json:"archive"`
		Author          string        `json:"author"`
		ChecksumMD5     string        `json:"checksum_md5"`
		ChecksumSHA256  string        `json:"checksum_sha256"`
		Date            string        `json:"date"`
		Distribution    string        `json:"Distribution"`
		DownloadURL     string        `json:"download_url"`
		License         []string      `json:"license"`
		MainModule      string        `json:"main_module"`
		Maturity        Maturity      `json:"maturity"`
		Metadata        cm.Spec       `json:"metadata"`
		Name            string        `json:"name"`
		Provides        []string      `json:"provides"`
		Resources       cm.Resources  `json:"resources"`
		Status          ReleaseStatus `json:"status"`
		Version         version.JSON  `json:"version"`
		VersionNumified float64       `json:"version_numified"`
		Authorized      bool          `json:"authorized"`
		First           bool          `json:"first"`
		Deprecated      bool          `json:"deprecated"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*r = Release{
		Dependency:      v.Dependency,
		Abstract:        v.Abstract,
		Archive:         v.Archive,
		Author:          v.Author,
		Authorized:      v.Authorized,
		ChecksumMD5:     v.ChecksumMD5,
		ChecksumSHA256:  v.ChecksumSHA256,
		Date:            v.Date,
		Deprecated:      v.Deprecated,
		Distribution:    v.Distribution,
		DownloadURL:     v.DownloadURL,
		First:           v.First,
		License:         v.License,
		MainModule:      v.MainModule,
		Maturity:        v.Maturity,
		Metadata:        v.Metadata,
		Name:            v.Name,
		Provides:        v.Provides,
		Resources:       v.Resources,
		Status:          v.Status,
		Version:         v.Version,
		VersionNumified: v.VersionNumified,
	}

	return nil
}

func (r *Release) _type() Type {
	return TypeRelease
}
