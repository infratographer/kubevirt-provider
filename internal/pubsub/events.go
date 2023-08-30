package pubsub

import (
	"context"
	"errors"

	"go.infratographer.com/x/events"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

var (
	tracer = otel.Tracer("go.equinixmetal.net/infra9-metal-bridge")

	// ErrNoEventsProvider is returned when no known providers have been configured.
	ErrNoEventsProvider = errors.New("no events provider configured")
)

// Events is the subscriber client
type Events struct {
	ctx            context.Context
	ChangeChannels []<-chan events.Message[events.ChangeMessage]
	logger         *zap.SugaredLogger
	conn           events.Connection
	// svc                service.Service
	maxProcessAttempts uint64
}

// Shutdown drains and closes the connection, blocking until complete.
func (e *Events) Shutdown(ctx context.Context) error {
	return e.conn.Shutdown(ctx)
}

// EventsOption is a functional option for the Events
type EventsOption func(s *Events)

// WithLogger sets the logger for the Events
func WithLogger(l *zap.SugaredLogger) EventsOption {
	return func(s *Events) {
		s.logger = l
	}
}

// NewEvents creates a new Events
func NewEvents(ctx context.Context, config Config, opts ...EventsOption) (*Events, error) {
	s := &Events{
		ctx:                ctx,
		logger:             zap.NewNop().Sugar(),
		maxProcessAttempts: config.MaxProcessAttempts,
	}

	for _, opt := range opts {
		opt(s)
	}

	events, err := events.NewConnection(config.Config, events.WithLogger(s.logger))
	if err != nil {
		return nil, err
	}

	s.conn = events

	return s, nil
}

// // SetService sets the service events will call.
// func (e *Events) SetService(service service.Service) {
// 	e.svc = service
// }
