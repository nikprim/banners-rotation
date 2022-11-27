package app

import "errors"

var (
	ErrNoOneBannerFoundForSlot   = errors.New("не найден ни 1 баннер для указанного слота")
	ErrBannerAlreadyLinkedToSlot = errors.New("баннер уже привязан к этому слоту")
	ErrBannerNotFound            = errors.New("баннер не найден")
	ErrSlotNotFound              = errors.New("слот не найден")
	ErrSocialGroupNotFound       = errors.New("социальная группа не найдена")
)
