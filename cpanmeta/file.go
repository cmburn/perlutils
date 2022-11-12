package cpanmeta

import "github.com/cmburn/perlutils/version"

type File struct {
	// File is the name of the file.
	File string `json:"file"`

	// Version is the version of the file.
	Version version.JSON `json:"version"`
}
