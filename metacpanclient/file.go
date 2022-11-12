package metacpanclient

import (
	"errors"
	"io"
	"strings"

	// local
	pui "github.com/cmburn/perlutils/internal"
	"github.com/cmburn/perlutils/version"

	// external
	"github.com/relvacode/iso8601"
)

type FileInfo struct {
	// Module is a list of modules indexed from within the file.
	Module []struct {
		// Name is the name of the module.
		Name string `json:"name"`

		// Indexed indicates whether the module is indexed by MetaCPAN.
		Indexed bool `json:"indexed"`

		// Authorized indicates whether the module is an authorized
		// upload.
		Authorized bool `json:"authorized"`

		// Version is the module's version.
		Version version.JSON `json:"version"`

		// VersionNumified is the version rendered as a float64
		VersionNumified float64 `json:"version_numified"`

		// AssociatedPod is the path to the POD documentation for the
		// given module.
		AssociatedPod string `json:"associated_pod"`
	} `json:"module"`

	// Stat contains info about the file similar to os.FileStat
	Stat Stat `json:"stat"`

	// Abstract contains the abstract portion of the name from the pod
	// documentation.
	Abstract string `json:"abstract"`

	// Author contains the author's Pause ID.
	Author string `json:"author"`

	// Date indicates when the file was uploaded.
	Date iso8601.Time `json:"date"`

	// Description contains the info from the DESCRIPTION field in a file's
	// documentation, if it exists.
	Description string `json:"description"`

	// Distribution is the name of the distribution that contains the file.
	Distribution string `json:"distribution"`

	// Documentation is the name of the module contained by this File.
	Documentation string `json:"documentation"`

	// DownloadURL is the url for the distribution archive for this file.
	DownloadURL string `json:"download_url"`

	// ID is the MetaCPAN ID for this file.
	ID string `json:"id"`

	// Level is the depth of the file from the archive root.
	Level int `json:"level"`

	// Maturity is the release maturity, i.e. release or developer
	Maturity Maturity `json:"maturity"`

	// Mime is the file's mime type.
	Mime string `json:"mime"`

	// Name is the file's name without any directory path
	Name string `json:"name"`

	// Path is the file path within the archive.
	Path string `json:"path"`

	// PodLines contains information about the documentation lines and
	// layout.
	PodLines [][]int `json:"pod_lines"`

	// Release is the release that contains this particular file.
	Release string `json:"release"`

	// SLOC is the number of lines of code in the file.
	SLOC int `json:"sloc"`

	// SLOP is the number of lines of POD in the file.
	SLOP int `json:"slop"`

	// Status is the release status for the file.
	Status ReleaseStatus `json:"status"`

	// Version is the version of the distribution that contains this file.
	Version version.JSON `json:"version"`

	// VersionNumified is the version of the distribution that contains
	// this file, in float format.
	VersionNumified float64 `json:"version_numified"`

	// Authorized indicates whether this file is from an authorized upload.
	Authorized bool `json:"authorized"`

	// Binary indicates whether this file has binary data.
	Binary bool `json:"binary"`

	// Deprecated indicates whether this file is deprecated.
	Deprecated bool `json:"deprecated"`

	// Directory indicates whether this file is a directory.
	Directory bool `json:"directory"`

	// Indexed indicates whether the content of the file is indexed.
	Indexed bool `json:"indexed"`
}

// File is an object specifying info about a file on MetaCPAN.
type File struct {
	FileInfo
	hasUA
	pod *Pod
	src string
}

func (f *File) MetaCPANURL() string {
	const baseURL = MetaCPANURL + "/source/"
	sb := strings.Builder{}
	sb.WriteString(baseURL)
	sb.WriteString(f.Author)
	sb.WriteRune('/')
	sb.WriteString(f.Release)
	sb.WriteRune('/')
	sb.WriteString(f.Path)
	return sb.String()
}

func (f *File) Pod(kind PodKind) (string, error) {
	if f.pod == nil {
		if f.mc == nil {
			return "", ErrNilClient
		}
		podName := f.Distribution
		podName = strings.ReplaceAll(podName, "-", "::")
		var err error
		if f.pod, err = NewPod(podName, "", f.mc); err != nil {
			return "", err
		}
	}
	return f.pod.get(kind)
}

func (f *File) Source() (string, error) {
	if f.src != "" {
		return f.src, nil
	}
	if f.mc == nil {
		return "", ErrNilClient
	}
	sb := strings.Builder{}
	sb.WriteString("/source/")
	sb.WriteString(f.Author)
	sb.WriteRune('/')
	sb.WriteString(f.Release)
	sb.WriteRune('/')
	sb.WriteString(f.Path)
	rc, err := f.mc.request("", sb.String(), nil)
	defer pui.CloseBody(rc)
	if err != nil {
		return "", err
	}
	buf, err := io.ReadAll(rc)
	if err != nil {
		return "", err
	}
	f.src = string(buf)
	return f.src, nil
}

func (f *File) _type() Type {
	return TypeFile
}

var (
	ErrInvalidPodKind = errors.New("invalid pod kind")
)
