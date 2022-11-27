package models

import (
	"github.com/google/uuid"
)

type Banner struct {
	GUID uuid.UUID
	Name string
}
