package app

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/nikprim/banners-rotation/internal/mab"
	rmqModels "github.com/nikprim/banners-rotation/internal/rmq/models"
	storageModels "github.com/nikprim/banners-rotation/internal/storage/models"
)

func (a *app) GetBanner(
	ctx context.Context,
	slotGUID,
	socialGroupGUID *uuid.UUID) (resultBanner *storageModels.Banner, err error) {
	defer func() {
		if err != nil {
			return
		}

		a.sendEventToQueue(ctx, &rmqModels.Event{
			Type:            rmqModels.EventTypeShow,
			BannerGUID:      resultBanner.GUID,
			SlotGUID:        *slotGUID,
			SocialGroupGUID: *socialGroupGUID,
			Datetime:        time.Now(),
		})
	}()

	err = a.checkSlotAndSocialGroupExists(ctx, slotGUID, socialGroupGUID)
	if err != nil {
		return nil, err
	}

	stats, err := a.storage.FindStatsBySlotAndSocialGroup(ctx, slotGUID, socialGroupGUID)
	if err != nil {
		return nil, err
	}

	bannersGUIDs, err := a.storage.FindBannersInSlot(ctx, slotGUID)
	if err != nil {
		return nil, err
	}

	if len(bannersGUIDs) == 0 {
		return nil, ErrNoOneBannerFoundForSlot
	}

	var (
		statsWithLink []*storageModels.Stat
		tries         int64
	)

	// учитываем статистику только для связанных баннеров и слотов
nextBanner:
	for _, bannerGUID := range bannersGUIDs {
		for _, stat := range stats {
			if stat.BannerGUID == *bannerGUID {
				// если статистика уже есть, то берём её
				statsWithLink = append(statsWithLink, stat)
				tries += int64(stat.Shows)

				continue nextBanner
			}
		}

		// если статистики ещё нет, то создаём новую
		statsWithLink = append(statsWithLink, &storageModels.Stat{
			BannerGUID:      *bannerGUID,
			SlotGUID:        *slotGUID,
			SocialGroupGUID: *socialGroupGUID,
			Shows:           0,
			Clicks:          0,
		})
	}

	resultStat := mab.UCB1(statsWithLink, tries)

	resultBanner, err = a.storage.FindBannerByGUID(ctx, &resultStat.BannerGUID)
	if err != nil {
		return nil, err
	}

	if resultBanner == nil {
		return nil, ErrBannerNotFound
	}

	if resultStat.GUID == uuid.Nil {
		resultStat.GUID = uuid.New()
		resultStat.Shows = 1

		err = a.storage.CreateStat(ctx, resultStat)
		if err != nil {
			return nil, err
		}

		return resultBanner, nil
	}

	err = a.storage.AddShowToStat(ctx, &resultStat.GUID)
	if err != nil {
		return nil, err
	}

	return resultBanner, nil
}

//nolint:dupl
func (a *app) checkSlotAndSocialGroupExists(ctx context.Context, slotGUID, socialGroupGUID *uuid.UUID) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var (
		slotErr        error
		socialGroupErr error
	)

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

	if slotErr != nil {
		return slotErr
	}

	if socialGroupErr != nil {
		return socialGroupErr
	}

	return nil
}
