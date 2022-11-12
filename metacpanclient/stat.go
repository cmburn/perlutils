package metacpanclient

import (
	"encoding/json"
	"io/fs"
	"time"

	// local
	pui "github.com/cmburn/perlutils/internal"
)

type Stat struct {
	ModTime time.Time   `json:"mtime"`
	Mode    fs.FileMode `json:"mode"`
	Size    int64       `json:"size"`
}

func (s *Stat) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]int64{
		"mtime": s.ModTime.Unix(),
		"mode":  int64(s.Mode),
		"size":  s.Size,
	})
}

func (s *Stat) UnmarshalJSON(data []byte) error {
	var v struct {
		ModTime pui.CoercedInt64  `json:"mtime"`
		Mode    pui.CoercedUInt32 `json:"mode"`
		Size    pui.CoercedInt    `json:"size"`
	}
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	s.ModTime = time.Unix(v.ModTime.Value, 0)
	s.Mode = fs.FileMode(v.Mode.Value)
	s.Size = int64(v.Size.Value)
	return nil
}
