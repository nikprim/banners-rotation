package apitest

import (
	"context"
	"time"

	"github.com/nikprim/banners-rotation/internal/server/pb"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type APISuiteActions struct {
	suite.Suite
	conn   *grpc.ClientConn
	client pb.BannerRotatorServiceClient
	ctx    context.Context

	amqpURI     string
	amqpQueue   string
	amqpConn    *amqp.Connection
	amqpChannel *amqp.Channel
}

func (s *APISuiteActions) Init(apiURL, amqpURI, amqpQueue string) {
	s.ctx = context.Background()

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*10)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		apiURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	s.Require().NoError(err)

	s.conn = conn
	s.client = pb.NewBannerRotatorServiceClient(conn)

	s.amqpURI = amqpURI
	s.amqpQueue = amqpQueue
}

func (s *APISuiteActions) End() {
	if err := s.conn.Close(); err != nil {
		s.Require().NoError(err)
	}
}

func (s *APISuiteActions) GetCountMessageInAMQP() int {
	if s.amqpConn == nil || s.amqpConn.IsClosed() {
		s.ConnectToAMQP()
	}

	queue, err := s.amqpChannel.QueueInspect(s.amqpQueue)
	s.Require().NoError(err)

	return queue.Messages
}

func (s *APISuiteActions) ConnectToAMQP() {
	var err error

	s.amqpConn, err = amqp.Dial(s.amqpURI)
	s.Require().NoError(err)

	s.amqpChannel, err = s.amqpConn.Channel()
	s.Require().NoError(err)
}
