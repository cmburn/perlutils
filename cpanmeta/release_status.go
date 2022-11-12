package cpanmeta

import (
	"fmt"

	// local
	"github.com/cmburn/perlutils/internal"
)

type ReleaseStatus int

const (
	ReleaseStatusStable ReleaseStatus = iota
	ReleaseStatusTesting
	ReleaseStatusUnstable
)

func (r *ReleaseStatus) String() string {
	switch *r {
	case ReleaseStatusStable:
		return "stable"
	case ReleaseStatusTesting:
		return "testing"
	case ReleaseStatusUnstable:
		return "unstable"
	default:
		return "undef"
	}
}

func NewReleaseStatus(str string) (ReleaseStatus, error) {
	switch str {
	case "stable":
		return ReleaseStatusStable, nil
	case "testing":
		return ReleaseStatusTesting, nil
	case "unstable":
		return ReleaseStatusUnstable, nil
	default:
		return ReleaseStatusStable,
			fmt.Errorf("invalid release status: %s", str)
	}
}

func (r *ReleaseStatus) MarshalJSON() ([]byte, error) {
	return internal.WrapEnumTypeJSON(r)
}

func (r *ReleaseStatus) UnmarshalJSON(data []byte) error {
	str := string(data)
	if len(str) < 2 {
		return fmt.Errorf("invalid release status: %s", str)
	}
	str = str[1 : len(str)-1]
	status, err := NewReleaseStatus(str)
	if err != nil {
		return err
	}
	*r = status
	return nil
}
