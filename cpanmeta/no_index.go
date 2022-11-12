package cpanmeta

type NoIndex struct {
	// File is a list of files that should not be indexed.
	File []string `json:"file"`

	// Directory is a list of directories that should not be indexed.
	Directory []string `json:"directory"`

	// Package is a list of packages that should not be indexed.
	Package []string `json:"package"`

	// Namespace is a list of namespaces that should not be indexed.
	Namespace []string `json:"namespace"`
}
