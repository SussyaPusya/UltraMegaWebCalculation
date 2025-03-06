package utils

import "github.com/google/uuid"

func IdGen() string {
	id := uuid.New()
	return id.String()
}
