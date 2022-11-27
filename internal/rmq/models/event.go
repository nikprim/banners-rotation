package models

import (
	"time"

	"github.com/google/uuid"
)

type EventType int

const (
	EventTypeShow EventType = iota
	EventTypeClick
)

type Event struct {
	Type            EventType `json:"type"`
	BannerGUID      uuid.UUID `json:"bannerGuid"`
	SlotGUID        uuid.UUID `json:"slotGuid"`
	SocialGroupGUID uuid.UUID `json:"socialGroupGuid"`
	Datetime        time.Time `json:"datetime"`
}
