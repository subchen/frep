package cli

import (
	"fmt"
	"time"
)

type timeValue struct {
	val *time.Time
}

func (v *timeValue) Set(value string) error {
	val, err := stringToTime(value)
	if err != nil {
		return err
	}

	*v.val = *val
	return nil
}

func (v *timeValue) String() string {
	return (*v.val).String()
}

var timePatterns = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05", // iso8601 without timezone
	time.RFC1123Z,
	time.RFC1123,
	"2006-01-02 15:04:05.999999999 -0700 MST",
	time.RFC822Z,
	time.RFC822,
	time.ANSIC,
	time.UnixDate,
	time.RubyDate,
	"2006-01-02 15:04:05Z07:00",
	"2006-01-02 15:04:05",
	"02 Jan 06 15:04 MST",
	"2006-01-02",
	"02 Jan 2006",
	"2006-01-02 15:04:05 -07:00",
	"2006-01-02 15:04:05 -0700",
}

func stringToTime(s string) (*time.Time, error) {
	for _, pattern := range timePatterns {
		if d, err := time.Parse(pattern, s); err == nil {
			return &d, nil
		}
	}
	return nil, fmt.Errorf("unable to parse date: %s", s)
}
