package misc

import (
	"errors"
	"time"
)

func ToUnixTime(t ...time.Time) int64 {

	if len(t) == 0 {
		return time.Now().Unix()
	}

	return t[0].Unix()
}

func FromUnixTime(unixTimestamp int64) (time.Time, error) {

	if unixTimestamp < 0 {
		return time.Time{}, errors.New("invalid Unix timestamp, must be non-negative")
	}

	return time.Unix(unixTimestamp, 0), nil
}

func FromUnixTime64(unixTimestamp int64) time.Time {
	return time.Unix(unixTimestamp, 0)
}

func FromUnixTime32(unixTimestamp int32) time.Time {
	return FromUnixTime64(int64(unixTimestamp))
}
