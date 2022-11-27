package storage

import (
	"context"
	"errors"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/google/uuid"
	pgx "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nikprim/banners-rotation/internal/app"
	"github.com/nikprim/banners-rotation/internal/storage/models"
)

type storage struct {
	conn *pgxpool.Pool
}

var _ app.Storage = (*storage)(nil)

func New(conn *pgxpool.Pool) app.Storage {
	return &storage{conn}
}

func (s *storage) AddBannerToSlot(ctx context.Context, bannerGUID, slotGUID *uuid.UUID) error {
	sql := `
		INSERT INTO banners_link_slots(guid, bannerGuid, slotGuid)
		VALUES ($1, $2, $3);
	`
	_, err := s.conn.Exec(ctx, sql, uuid.New(), bannerGUID, slotGUID)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) FindLinkBannerAndSlot(
	ctx context.Context,
	bannerGUID,
	slotGUID *uuid.UUID) (*models.LinkBannerAndSlot, error) {
	query := `
		SELECT guid, bannerGuid, slotGuid
		FROM banners_link_slots
		WHERE bannerGuid = $1 AND slotGuid = $2
	`

	var link models.LinkBannerAndSlot

	err := pgxscan.Get(ctx, s.conn, &link, query, bannerGUID, slotGUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &link, nil
}

func (s *storage) RemoveBannerFromSlot(ctx context.Context, bannerGUID, slotGUID *uuid.UUID) error {
	sql := `
		DELETE FROM banners_link_slots
		WHERE bannerGuid = $1 AND slotGuid = $2
	`
	_, err := s.conn.Exec(ctx, sql, bannerGUID, slotGUID)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) FindStatByParams(
	ctx context.Context,
	bannerGUID,
	slotGUID,
	socialGroupGUID *uuid.UUID) (*models.Stat, error) {
	query := `
		SELECT guid, bannerGuid, slotGuid, socialGroupGuid, shows, clicks
		FROM stats
		WHERE bannerGuid = $1 AND slotGuid = $2 AND socialGroupGuid = $3
	`

	var stat models.Stat

	err := pgxscan.Get(ctx, s.conn, &stat, query, bannerGUID, slotGUID, socialGroupGUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &stat, nil
}

func (s *storage) CreateStat(ctx context.Context, val *models.Stat) error {
	sql := `
		INSERT INTO stats(guid, bannerGuid, slotGuid, socialGroupGuid, shows, clicks)
		VALUES ($1, $2, $3, $4, $5, $6);
	`
	_, err := s.conn.Exec(ctx, sql, val.GUID, val.BannerGUID, val.SlotGUID, val.SocialGroupGUID, val.Shows, val.Clicks)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) AddClickToStat(ctx context.Context, statGUID *uuid.UUID) error {
	sql := `
		UPDATE stats
		SET clicks = clicks + 1
		WHERE guid = $1
	`

	_, err := s.conn.Exec(ctx, sql, statGUID)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) AddShowToStat(ctx context.Context, statGUID *uuid.UUID) error {
	sql := `
		UPDATE stats
		SET shows = shows + 1
		WHERE guid = $1
	`

	_, err := s.conn.Exec(ctx, sql, statGUID)
	if err != nil {
		return err
	}

	return nil
}

func (s *storage) FindStatsBySlotAndSocialGroup(
	ctx context.Context,
	slotGUID,
	statGUID *uuid.UUID) ([]*models.Stat, error) {
	sql := `
		SELECT guid, bannerGuid, slotGuid, socialGroupGuid, shows, clicks
		FROM stats
		WHERE slotGuid = $1 AND socialGroupGuid = $2
	`

	var stats []*models.Stat

	err := pgxscan.Select(ctx, s.conn, &stats, sql, slotGUID, statGUID)
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (s *storage) FindBannersInSlot(ctx context.Context, slotGUID *uuid.UUID) ([]*uuid.UUID, error) {
	sql := `
		SELECT bannerGuid
		FROM banners_link_slots
		WHERE slotGuid = $1
	`

	var bannersGUIDs []*uuid.UUID

	err := pgxscan.Select(ctx, s.conn, &bannersGUIDs, sql, slotGUID)
	if err != nil {
		return nil, err
	}

	return bannersGUIDs, nil
}

func (s *storage) FindBannerByGUID(ctx context.Context, bannerGUID *uuid.UUID) (*models.Banner, error) {
	query := `
		SELECT guid, name
		FROM banners
		WHERE guid = $1
	`

	var banner models.Banner

	err := pgxscan.Get(ctx, s.conn, &banner, query, bannerGUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &banner, nil
}

func (s *storage) FindSlotByGUID(ctx context.Context, slotGUID *uuid.UUID) (*models.Slot, error) {
	query := `
		SELECT guid, name
		FROM slots
		WHERE guid = $1
	`

	var slot models.Slot

	err := pgxscan.Get(ctx, s.conn, &slot, query, slotGUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &slot, nil
}

func (s *storage) FindSocialGroupByGUID(ctx context.Context, socialGroupGUID *uuid.UUID) (*models.SocialGroup, error) {
	query := `
		SELECT guid, name
		FROM social_groups
		WHERE guid = $1
	`

	var socialGroup models.SocialGroup

	err := pgxscan.Get(ctx, s.conn, &socialGroup, query, socialGroupGUID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &socialGroup, nil
}
