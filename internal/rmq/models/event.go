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
	BannerGUID      uuid.UUID `json:"banner_guid"`
	SlotGUID        uuid.UUID `json:"slot_guid"`
	SocialGroupGUID uuid.UUID `json:"social_group_guid"`
	Datetime        time.Time `json:"datetime"`
}
