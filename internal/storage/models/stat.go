package models

import (
	"github.com/google/uuid"
)

type Stat struct {
	GUID            uuid.UUID
	BannerGUID      uuid.UUID
	SlotGUID        uuid.UUID
	SocialGroupGUID uuid.UUID
	Shows           int
	Clicks          int
}
