package app

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	rmqModels "github.com/nikprim/banners-rotation/internal/rmq/models"
	storageModels "github.com/nikprim/banners-rotation/internal/storage/models"
)

func (a *app) AddClick(ctx context.Context, bannerGUID, slotGUID, socialGroupGUID *uuid.UUID) (err error) {
	defer func() {
		if err != nil {
			return
		}

		a.sendEventToQueue(ctx, &rmqModels.Event{
			Type:            rmqModels.EventTypeClick,
			BannerGUID:      *bannerGUID,
			SlotGUID:        *slotGUID,
			SocialGroupGUID: *socialGroupGUID,
			Datetime:        time.Now(),
		})
	}()

	err = a.checkBannerAndSlotAndSocialGroupExists(ctx, bannerGUID, slotGUID, socialGroupGUID)
	if err != nil {
		return err
	}

	stat, err := a.storage.FindStatByParams(ctx, bannerGUID, slotGUID, socialGroupGUID)
	if err != nil {
		return err
	}

	if stat == nil {
		err = a.storage.CreateStat(ctx, &storageModels.Stat{
			GUID:            uuid.New(),
			BannerGUID:      *bannerGUID,
			SlotGUID:        *slotGUID,
			SocialGroupGUID: *socialGroupGUID,
			Shows:           0,
			Clicks:          1,
		})
		if err != nil {
			return err
		}

		return nil
	}

	return a.storage.AddClickToStat(ctx, &stat.GUID)
}

func (a *app) checkBannerAndSlotAndSocialGroupExists(
	ctx context.Context,
	bannerGUID,
	slotGUID,
	socialGroupGUID *uuid.UUID) error {
	wg := sync.WaitGroup{}
	wg.Add(3)

	var (
		bannerErr      error
		slotErr        error
		socialGroupErr error
	)

	go func() {
		defer wg.Done()

		var banner *storageModels.Banner

		banner, bannerErr = a.storage.FindBannerByGUID(ctx, bannerGUID)
		if bannerErr != nil {
			return
		}

		if banner == nil {
			bannerErr = ErrBannerNotFound
		}
	}()

	go func() {
		defer wg.Done()

		var slot *storageModels.Slot

		slot, slotErr = a.storage.FindSlotByGUID(ctx, slotGUID)
		if slotErr != nil {
			return
		}

		if slot == nil {
			slotErr = ErrSlotNotFound
		}
	}()

	go func() {
		defer wg.Done()

		var socialGroup *storageModels.SocialGroup

		socialGroup, socialGroupErr = a.storage.FindSocialGroupByGUID(ctx, socialGroupGUID)
		if socialGroupErr != nil {
			return
		}

		if socialGroup == nil {
			socialGroupErr = ErrSocialGroupNotFound
		}
	}()

	wg.Wait()

	if bannerErr != nil {
		return bannerErr
	}

	if slotErr != nil {
		return slotErr
	}

	if socialGroupErr != nil {
		return socialGroupErr
	}

	return nil
}
