package app

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/nikprim/banners-rotation/internal/rmq"
	rmqModels "github.com/nikprim/banners-rotation/internal/rmq/models"
	storageModels "github.com/nikprim/banners-rotation/internal/storage/models"
	"github.com/rs/zerolog/log"
)

type Application interface {
	AddBannerToSlot(ctx context.Context, bannerGUID, slotGUID *uuid.UUID) error
	RemoveBannerFromSlot(ctx context.Context, bannerGUID, slotGUID *uuid.UUID) error
	AddClick(ctx context.Context, bannerGUID, slotGUID, socialGroup *uuid.UUID) error
	GetBanner(ctx context.Context, slotGUID, socialGroup *uuid.UUID) (*storageModels.Banner, error)
}

type app struct {
	storage  Storage
	producer *rmq.Producer
}

type Storage interface {
	AddBannerToSlot(ctx context.Context, bannerGUID, slotGUID *uuid.UUID) error
	FindLinkBannerAndSlot(ctx context.Context, bannerGUID, slotGUID *uuid.UUID) (*storageModels.LinkBannerAndSlot, error)
	RemoveBannerFromSlot(ctx context.Context, bannerGUID, slotGUID *uuid.UUID) error
	FindStatByParams(ctx context.Context, bannerGUID, slotGUID, socialGroupGUID *uuid.UUID) (*storageModels.Stat, error)
	CreateStat(ctx context.Context, val *storageModels.Stat) error
	AddClickToStat(ctx context.Context, statGUID *uuid.UUID) error
	AddShowToStat(ctx context.Context, statGUID *uuid.UUID) error
	FindStatsBySlotAndSocialGroup(ctx context.Context, slotGUID, statGUID *uuid.UUID) ([]*storageModels.Stat, error)
	FindBannersInSlot(ctx context.Context, slotGUID *uuid.UUID) ([]*uuid.UUID, error)
	FindBannerByGUID(ctx context.Context, bannerGUID *uuid.UUID) (*storageModels.Banner, error)
	FindSlotByGUID(ctx context.Context, slotGUID *uuid.UUID) (*storageModels.Slot, error)
	FindSocialGroupByGUID(ctx context.Context, socialGroupGUID *uuid.UUID) (*storageModels.SocialGroup, error)
}

func New(storage Storage, producer *rmq.Producer) Application {
	return &app{
		storage:  storage,
		producer: producer,
	}
}

func (a *app) sendEventToQueue(ctx context.Context, message *rmqModels.Event) {
	l := log.With().Fields(map[string]interface{}{
		"type":            message.Type,
		"bannerGuid":      message.BannerGUID.String(),
		"slotGuid":        message.SlotGUID.String(),
		"socialGroupGuid": message.SocialGroupGUID.String(),
		"datetime":        message.Datetime.Format(time.RFC3339),
	}).Logger()

	body, err := json.Marshal(message)
	if err != nil {
		l.Error().Err(err).Msg("failed to send event to queue")

		return
	}

	err = a.producer.Publish(ctx, body)
	if err != nil {
		l.Error().Err(err).Msg("failed to send event to queue")

		return
	}

	l.Info().Msg("event sent to queue")
}
