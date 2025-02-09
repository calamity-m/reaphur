package util

import "github.com/google/uuid"

func ParseUUIDRegardless(u string) uuid.UUID {
	parsed, err := uuid.Parse(u)
	if err != nil {
		return uuid.Nil
	}

	return parsed
}
