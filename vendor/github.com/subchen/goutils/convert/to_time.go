package convert

import (
	"fmt"
	"time"
)

func AsTime(value interface{}) time.Time {
	v, _ := toTime(value)
	return v
}

func ToTime(value interface{}) (time.Time, error) {
	return toTime(value)
}

func toTime(value interface{}) (time.Time, error) {
	switch v := value.(type) {
	case time.Time:
		return v, nil
	case string:
		return stringToTime(v)
	default:
		return time.Time{}, fmt.Errorf("unable to cast %T to time.Time", value)
	}
}

var timePatterns = []string{
	time.RFC3339,
	time.RFC3339Nano,
	"2006-01-02T15:04:05", // iso8601 without timezone
	time.RFC1123Z,
	time.RFC1123,
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

func stringToTime(s string) (time.Time, error) {
	for _, pattern := range timePatterns {
		if d, err := time.Parse(pattern, s); err == nil {
			return d, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", s)
}
