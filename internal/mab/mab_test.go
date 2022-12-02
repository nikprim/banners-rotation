package mab_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nikprim/banners-rotation/internal/mab"
	storageModels "github.com/nikprim/banners-rotation/internal/storage/models"
	"github.com/stretchr/testify/require"
)

func TestUCB1(t *testing.T) {
	stats := make([]*storageModels.Stat, 0, 10)

	for i := 0; i < cap(stats); i++ {
		stats = append(stats, &storageModels.Stat{
			GUID:            uuid.New(),
			BannerGUID:      uuid.New(),
			SlotGUID:        uuid.New(),
			SocialGroupGUID: uuid.New(),
		})
	}

	t.Run("Баннер показывается чаще остальных", func(t *testing.T) {
		popularBannerGUID := uuid.New()
		stats = append(stats, &storageModels.Stat{
			BannerGUID:      popularBannerGUID,
			SlotGUID:        uuid.New(),
			SocialGroupGUID: uuid.New(),
			Clicks:          10,
		})

		var (
			maxShows                int
			resultPopularBannerGUID uuid.UUID
		)

		for i := len(stats) + 1; i < len(stats)*1000; i++ {
			stat := mab.UCB1(stats, int64(i))
			stat.Shows++

			if stat.Shows > maxShows {
				maxShows = stat.Shows
				resultPopularBannerGUID = stat.BannerGUID
			}
		}

		require.Equal(t, popularBannerGUID, resultPopularBannerGUID)
	})

	t.Run("Показаны все баннеры", func(t *testing.T) {
		for i := 1; i <= len(stats); i++ {
			stat := mab.UCB1(stats, int64(i))
			stat.Shows++
		}

		for _, stat := range stats {
			require.NotEqual(t, 0, stat.Shows)
		}
	})
}
