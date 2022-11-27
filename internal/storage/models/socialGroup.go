package models

import (
	"github.com/google/uuid"
)

type SocialGroup struct {
	GUID uuid.UUID
	Name string
}
