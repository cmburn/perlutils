package metacpanclient

import (
	"github.com/cmburn/perlutils/version"
	"github.com/relvacode/iso8601"
)

// DownloadURL contains information about a download URL for a given release.
type DownloadURL struct {
	// ChecksumSHA256 is the SHA256 checksum of the release.
	ChecksumSHA256 string `json:"checksum_sha256"`

	// ChecksumMD5 is the MD5 checksum of the release.
	ChecksumMD5 string `json:"checksum_md5"`

	// Date is the release date for a given distribution
	Date iso8601.Time `json:"date"`

	// DownloadURL is the URL to download the release archive.
	DownloadURL string `json:"download_url"`

	// Status is the release status for the given distribution.
	Status ReleaseStatus `json:"status"`

	// Version is the specific version this DownloadURL pertains to.
	Version version.JSON `json:"version"`
}
