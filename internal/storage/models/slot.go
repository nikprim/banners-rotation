package models

import (
	"github.com/google/uuid"
)

type Slot struct {
	GUID uuid.UUID
	Name string
}
