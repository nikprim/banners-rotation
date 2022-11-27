package service

import (
	"context"

	"github.com/google/uuid"
	internalApp "github.com/nikprim/banners-rotation/internal/app"
	"github.com/nikprim/banners-rotation/internal/server/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

type bannerRotatorService struct {
	pb.UnimplementedBannerRotatorServiceServer
	app internalApp.Application
}

func NewBannerRotatorService(app internalApp.Application) pb.BannerRotatorServiceServer {
	return &bannerRotatorService{
		app: app,
	}
}

func (e *bannerRotatorService) AddBannerToSlot(
	ctx context.Context,
	request *pb.BannerAndSlotRequest) (*emptypb.Empty, error) {
	bannerGUID, err := uuid.Parse(request.GetBannerGuid())
	if err != nil {
		return nil, err
	}

	slotGUID, err := uuid.Parse(request.GetSlotGuid())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, e.app.AddBannerToSlot(ctx, &bannerGUID, &slotGUID)
}

func (e *bannerRotatorService) RemoveBannerFromSlot(
	ctx context.Context,
	request *pb.BannerAndSlotRequest) (*emptypb.Empty, error) {
	bannerGUID, err := uuid.Parse(request.GetBannerGuid())
	if err != nil {
		return nil, err
	}

	slotGUID, err := uuid.Parse(request.GetSlotGuid())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, e.app.RemoveBannerFromSlot(ctx, &bannerGUID, &slotGUID)
}

func (e *bannerRotatorService) AddClick(ctx context.Context, request *pb.AddClickRequest) (*emptypb.Empty, error) {
	bannerGUID, err := uuid.Parse(request.GetBannerGuid())
	if err != nil {
		return nil, err
	}

	slotGUID, err := uuid.Parse(request.GetSlotGuid())
	if err != nil {
		return nil, err
	}

	socialGroupGUID, err := uuid.Parse(request.GetSocialGroupGuid())
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, e.app.AddClick(ctx, &bannerGUID, &slotGUID, &socialGroupGUID)
}

func (e *bannerRotatorService) GetBanner(
	ctx context.Context,
	request *pb.SlotAndSocialGroupRequest) (*pb.Banner, error) {
	slotGUID, err := uuid.Parse(request.GetSlotGuid())
	if err != nil {
		return nil, err
	}

	socialGroupGUID, err := uuid.Parse(request.GetSocialGroupGuid())
	if err != nil {
		return nil, err
	}

	banner, err := e.app.GetBanner(ctx, &slotGUID, &socialGroupGUID)
	if err != nil {
		return nil, err
	}

	return &pb.Banner{
		Guid: banner.GUID.String(),
		Name: banner.Name,
	}, nil
}
