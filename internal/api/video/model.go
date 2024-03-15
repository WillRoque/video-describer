package video

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
)

type DescribeRequest struct {
	VideoID string   `json:"videoID"`
	Prompt  string   `json:"prompt"`
	Start   Duration `json:"start"`
	End     Duration `json:"end"`
}

func (dr *DescribeRequest) Validate() string {
	var v []string
	if dr.VideoID == "" {
		v = append(v, "missing video ID")
	}
	if (dr.Start.Duration != 0 && dr.End.Duration != 0) && dr.Start.Duration >= dr.End.Duration {
		v = append(v, "start must be before end")
	}
	return strings.Join(v, ",")
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return fmt.Errorf("error unmarshalling JSON: %w", err)
	}
	switch value := v.(type) {
	case float64:
		d.Duration = time.Duration(value)
		return nil
	case string:
		var err error
		d.Duration, err = time.ParseDuration(value)
		if err != nil {
			return fmt.Errorf("error parsing duration string: %w", err)
		}
		return nil
	default:
		return errors.New("invalid duration type")
	}
}
