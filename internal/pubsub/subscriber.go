package pubsub

import (
	"errors"
	"sync"
	"time"

	"go.infratographer.com/x/events"
)

const defaultNakDelay = 10 * time.Second

// ErrServiceRequired is returned when Events.Listen is called before Events.SetService
var ErrServiceRequired = errors.New("event service required")

// Subscribe subscribes to a nats subject
func (e *Events) Subscribe(topic string) error {
	msgChan, err := e.conn.SubscribeChanges(e.ctx, topic)
	if err != nil {
		return err
	}

	e.ChangeChannels = append(e.ChangeChannels, msgChan)

	e.logger.Infof("Subscribing to topic %s", topic)

	return nil
}

// Listen start listening for messages on registered subjects and calls the registered message handler
func (e *Events) Listen() error {
	// if e.svc == nil {
	// 	return ErrServiceRequired
	// }

	wg := &sync.WaitGroup{}

	// goroutine for each change channel
	for _, ch := range e.ChangeChannels {
		wg.Add(1)

		go e.listen(ch, wg)
	}

	wg.Wait()

	return nil
}

// listen listens for messages on a channel and calls the registered message handler
func (e *Events) listen(messages <-chan events.Message[events.ChangeMessage], wg *sync.WaitGroup) {
	defer wg.Done()

	for msg := range messages {
		elogger := e.logger.With(
			"event.message.id", msg.ID(),
			"event.message.topic", msg.Topic(),
			"event.message.timestamp", msg.Timestamp(),
			"event.message.deliveries", msg.Deliveries(),
		)

		if err := e.processEvent(msg); err != nil {
			elogger.Errorw("failed to process msg", "error", err)

			if e.maxProcessAttempts != 0 && msg.Deliveries()+1 > e.maxProcessAttempts {
				elogger.Warnw("terminating event, too many attempts")

				if termErr := msg.Term(); termErr != nil {
					elogger.Warnw("error occurred while terminating event")
				}
			} else if nakErr := msg.Nak(defaultNakDelay); nakErr != nil {
				elogger.Warnw("error occurred while naking", "error", nakErr)
			}
		} else if ackErr := msg.Ack(); ackErr != nil {
			elogger.Warnw("error occurred while acking", "error", ackErr)
		}
	}
}

// processEvent event message handler
func (e *Events) processEvent(msg events.Message[events.ChangeMessage]) error {
	elogger := e.logger.With(
		"event.message.id", msg.ID(),
		"event.message.timestamp", msg.Timestamp(),
		"event.message.deliveries", msg.Deliveries(),
	)

	if msg.Error() != nil {
		elogger.Errorw("message contains error:", "error", msg.Error())

		return msg.Error()
	}

	return nil
}

/*
	changeMsg := msg.Message()

	ctx := changeMsg.GetTraceContext(context.Background())

	ctx, span := tracer.Start(ctx, "pubsub.receive",
		trace.WithAttributes(
			attribute.String("pubsub.message.id", msg.ID()),
			attribute.String("pubsub.subject", msg.Message().SubjectID.String()),
		),
	)

	defer span.End()

	elogger = elogger.With(
		"event.resource.id", changeMsg.SubjectID.String(),
		"event.type", changeMsg.EventType,
	)

	elogger.Debugw("received message")

	var err error

	switch events.ChangeType(msg.Message().EventType) {
	case events.CreateChangeType:
		err = e.handleTouchEvent(ctx, msg)
	case events.UpdateChangeType:
		err = e.handleTouchEvent(ctx, msg)
	case events.DeleteChangeType:
		err = e.handleDeleteEvent(ctx, msg)
	default:
		elogger.Warn("ignoring msg, not a create, update or delete event")
	}

	if err != nil {
		return err
	}

	return nil
}


func (e *Events) handleDeleteEvent(ctx context.Context, msg events.Message[events.ChangeMessage]) error {
	elogger := e.logger.With(
		"event.message.id", msg.ID(),
		"event.message.topic", msg.Topic(),
		"event.message.timestamp", msg.Timestamp(),
		"event.message.deliveries", msg.Deliveries(),
		"event.resource.id", msg.Message().SubjectID.String(),
		"event.type", msg.Message().EventType,
	)

	if e.svc.IsOrganizationID(msg.Message().SubjectID) {
		if err := e.svc.DeleteOrganization(ctx, msg.Message().SubjectID); err != nil {
			// TODO: only return errors on retryable errors
			return err
		}

		return nil
	}

	if e.svc.IsProjectID(msg.Message().SubjectID) {
		if err := e.svc.DeleteProject(ctx, msg.Message().SubjectID); err != nil {
			// TODO: only return errors on retryable errors
			return err
		}

		return nil
	}

	if e.svc.IsUser(msg.Message().SubjectID) {
		userUUID := msg.Message().SubjectID.String()[gidx.PrefixPartLength+1:]

		subjID, err := models.GenerateSubjectID(models.IdentityPrefixUser, models.MetalUserIssuer, models.MetalUserIssuerIDPrefix+userUUID)
		if err != nil {
			elogger.Errorw("failed to convert user id to identity id", "error", err)

			return nil
		}

		if err := e.svc.UnassignUser(ctx, subjID, msg.Message().AdditionalSubjectIDs...); err != nil {
			// TODO: only return errors on retryable errors
			return err
		}

		return nil
	}

	elogger.Warnw("unknown subject id")

	return nil
}
*/
