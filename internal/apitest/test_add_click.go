package apitest

import (
	"time"

	"github.com/nikprim/banners-rotation/internal/app"
	"github.com/nikprim/banners-rotation/internal/server/pb"
)

func (s *APISuite) TestAddClickSuccess() {
	countMessage := s.GetCountMessageInAMQP()

	_, err := s.client.AddClick(s.ctx, &pb.AddClickRequest{
		BannerGuid:      BannerGUID1,
		SlotGuid:        SlotGUID1,
		SocialGroupGuid: SocialGroupGUID1,
	})
	s.Require().NoError(err)
	time.Sleep(time.Millisecond * 100) // в CI проверка количества сообщений падает
	s.Require().Equal(countMessage+1, s.GetCountMessageInAMQP())

	_, err = s.client.AddClick(s.ctx, &pb.AddClickRequest{
		BannerGuid:      BannerGUID1,
		SlotGuid:        SlotGUID1,
		SocialGroupGuid: SocialGroupGUID1,
	})
	s.Require().NoError(err)
	time.Sleep(time.Millisecond * 100) // в CI проверка количества сообщений падает
	s.Require().Equal(countMessage+2, s.GetCountMessageInAMQP())
}

func (s *APISuite) TestAddClickErrors() {
	countMessage := s.GetCountMessageInAMQP()

	_, err := s.client.AddClick(s.ctx, &pb.AddClickRequest{
		BannerGuid:      BannerGUIDNotFound,
		SlotGuid:        SlotGUID1,
		SocialGroupGuid: SocialGroupGUID1,
	})
	s.Require().ErrorContains(err, app.ErrBannerNotFound.Error())
	s.Require().Equal(countMessage, s.GetCountMessageInAMQP())

	_, err = s.client.AddClick(s.ctx, &pb.AddClickRequest{
		BannerGuid:      BannerGUID1,
		SlotGuid:        SlotGUIDNotFound,
		SocialGroupGuid: SocialGroupGUID1,
	})
	s.Require().ErrorContains(err, app.ErrSlotNotFound.Error())
	s.Require().Equal(countMessage, s.GetCountMessageInAMQP())

	_, err = s.client.AddClick(s.ctx, &pb.AddClickRequest{
		BannerGuid:      BannerGUID1,
		SlotGuid:        SlotGUID1,
		SocialGroupGuid: SlotGUIDNotFound,
	})
	s.Require().ErrorContains(err, app.ErrSocialGroupNotFound.Error())
	s.Require().Equal(countMessage, s.GetCountMessageInAMQP())
}
