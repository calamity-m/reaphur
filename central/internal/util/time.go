package util

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func ParseProtoTimestamp(timestamp *timestamppb.Timestamp) time.Time {
	if timestamp == nil {
		return time.Time{}
	}

	if timestamp.IsValid() {
		if err := timestamp.CheckValid(); err == nil {
			return timestamp.AsTime()
		}
	}

	return time.Time{}
}
