package models

import (
	"github.com/google/uuid"
)

type LinkBannerAndSlot struct {
	GUID       uuid.UUID
	BannerGUID uuid.UUID
	SlotGUID   uuid.UUID
}
