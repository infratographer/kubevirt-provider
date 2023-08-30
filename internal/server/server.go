package server

import (
	"context"
	"errors"
	"net/http"

	"go.infratographer.com/kubevirt-provider/internal/pubsub"
	"go.infratographer.com/x/echox"
	"go.infratographer.com/x/events"
	"go.uber.org/zap"

	"go.infratographer.com/ipam-api/pkg/ipamclient"
)

// Server holds options for server connectivity and settings
type Server struct {
	APIClient    *http.Client
	IPAMClient   *ipamclient.Client
	Context      context.Context
	Debug        bool
	Echo         *echox.Server
	IPBlock      string
	Locations    []string
	Logger       *zap.SugaredLogger
	Events       events.Config
	ChangeTopics []string
	Publisher    events.Publisher

	events pubsub.Events
}

// Run will start the server queue connections and healthcheck endpoints
func (s *Server) Run(ctx context.Context) error {
	go func() {
		if err := s.Echo.Run(); err != nil {
			s.Logger.Error("unable to start healthcheck server", zap.Error(err))
		}
	}()

	s.Logger.Infow("starting subscribers")

	if err := s.ConfigureSubscribers(); err != nil {
		s.Logger.Errorw("unable to configure subscribers", "error", err)
		return err
	}

	for _, ch := range s.events.ChangeChannels {
		go s.ProcessChange(ch)
	}

	return nil
}

func (s *Server) ConfigureSubscribers() error {
	var err error

	for _, topic := range s.ChangeTopics {
		s.Logger.Debugw("subscribing to topic", "topic", topic)

		conn, err := events.NewConnection(s.Events, events.WithLogger(s.Logger))
		if err != nil {
			s.Logger.Errorw("unable to create change subscriber", "error", err, "topic", topic)
			err = errors.Join(err, errSubscriberCreate)
		}

		changes, err := conn.SubscribeChanges(s.Context, topic)
		if err != nil {
			s.Logger.Errorw("unable to subscribe to change topic", "error", err, "topic", topic, "type", "change")
			err = errors.Join(err, errSubscriptionCreate)
		}

		s.events.ChangeChannels = append(s.events.ChangeChannels, changes)
	}

	return err
}
