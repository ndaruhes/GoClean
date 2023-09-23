package utils

import "time"

func GetStringPointerValue(key *string) string {
	if key != nil {
		return *key
	}

	return ""
}

func GetTimePointerValue(key *time.Time) time.Time {
	if key != nil {
		return *key
	}

	return time.Time{}
}
