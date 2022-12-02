package mab

import (
	"math"

	storageModels "github.com/nikprim/banners-rotation/internal/storage/models"
)

func UCB1(stats []*storageModels.Stat, tries int64) *storageModels.Stat {
	var (
		maxConfidence  float64
		rotationToShow *storageModels.Stat
	)

	for _, stat := range stats {
		meanClicks := float64(stat.Clicks) / float64(stat.Shows+1)

		confidence := meanClicks + math.Sqrt(2*math.Log(float64(tries+1))/float64(stat.Shows+1))

		if confidence >= maxConfidence {
			maxConfidence = confidence
			rotationToShow = stat
		}
	}

	return rotationToShow
}
