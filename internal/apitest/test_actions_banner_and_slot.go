package apitest

import (
	"github.com/nikprim/banners-rotation/internal/app"
	"github.com/nikprim/banners-rotation/internal/server/pb"
)

func (s *APISuite) TestActionsBannerAndSlotSuccess() {
	_, err := s.client.AddBannerToSlot(s.ctx, &pb.BannerAndSlotRequest{
		BannerGuid: LinkNotFound1.bannerGUID,
		SlotGuid:   LinkNotFound1.slotGUID,
	})
	s.Require().NoError(err)

	_, err = s.client.RemoveBannerFromSlot(s.ctx, &pb.BannerAndSlotRequest{
		BannerGuid: LinkNotFound1.bannerGUID,
		SlotGuid:   LinkNotFound1.slotGUID,
	})
	s.Require().NoError(err)
}

func (s *APISuite) TestAddBannerToSlotErrors() {
	_, err := s.client.AddBannerToSlot(s.ctx, &pb.BannerAndSlotRequest{
		BannerGuid: BannerGUIDNotFound,
		SlotGuid:   SlotGUID1,
	})
	s.Require().ErrorContains(err, app.ErrBannerNotFound.Error())

	_, err = s.client.AddBannerToSlot(s.ctx, &pb.BannerAndSlotRequest{
		BannerGuid: BannerGUID1,
		SlotGuid:   SlotGUIDNotFound,
	})
	s.Require().ErrorContains(err, app.ErrSlotNotFound.Error())

	_, err = s.client.AddBannerToSlot(s.ctx, &pb.BannerAndSlotRequest{
		BannerGuid: Link1.bannerGUID,
		SlotGuid:   Link1.slotGUID,
	})
	s.Require().ErrorContains(err, app.ErrBannerAlreadyLinkedToSlot.Error())
}

func (s *APISuite) TestRemoveBannerFromSlotErrors() {
	_, err := s.client.RemoveBannerFromSlot(s.ctx, &pb.BannerAndSlotRequest{
		BannerGuid: BannerGUIDNotFound,
		SlotGuid:   SlotGUID1,
	})
	s.Require().ErrorContains(err, app.ErrBannerNotFound.Error())

	_, err = s.client.RemoveBannerFromSlot(s.ctx, &pb.BannerAndSlotRequest{
		BannerGuid: BannerGUID1,
		SlotGuid:   SlotGUIDNotFound,
	})
	s.Require().ErrorContains(err, app.ErrSlotNotFound.Error())

	_, err = s.client.RemoveBannerFromSlot(s.ctx, &pb.BannerAndSlotRequest{
		BannerGuid: LinkNotFound1.bannerGUID,
		SlotGuid:   LinkNotFound1.slotGUID,
	})
	s.Require().NoError(err)
}
