package app

import (
	"context"
	"sync"

	"github.com/google/uuid"
	storageModels "github.com/nikprim/banners-rotation/internal/storage/models"
)

func (a *app) AddBannerToSlot(ctx context.Context, bannerGUID, slotGUID *uuid.UUID) error {
	err := a.checkBannerAndSlotExists(ctx, bannerGUID, slotGUID)
	if err != nil {
		return err
	}

	link, err := a.storage.FindLinkBannerAndSlot(ctx, bannerGUID, slotGUID)
	if err != nil {
		return err
	}

	if link != nil {
		return ErrBannerAlreadyLinkedToSlot
	}

	return a.storage.AddBannerToSlot(ctx, bannerGUID, slotGUID)
}

func (a *app) RemoveBannerFromSlot(ctx context.Context, bannerGUID, slotGUID *uuid.UUID) error {
	err := a.checkBannerAndSlotExists(ctx, bannerGUID, slotGUID)
	if err != nil {
		return err
	}

	return a.storage.RemoveBannerFromSlot(ctx, bannerGUID, slotGUID)
}

//nolint:dupl
func (a *app) checkBannerAndSlotExists(ctx context.Context, bannerGUID, slotGUID *uuid.UUID) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	var (
		bannerErr error
		slotErr   error
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

	wg.Wait()

	if bannerErr != nil {
		return bannerErr
	}

	if slotErr != nil {
		return slotErr
	}

	return nil
}
