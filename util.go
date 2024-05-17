package plex

import (
	"time"
)

// GetSeq msg seq ms
func GetSeq() int64 {
	return time.Now().UnixNano() / 1000000
}
