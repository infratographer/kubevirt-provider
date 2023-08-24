package pubsub

import (
	"context"

	"go.infratographer.com/x/events"
	"go.uber.org/multierr"
)

// PublishAuthRelationshipRequest publishes an auth relationship request message to the event handler, and waits for a response.
func (e *Events) PublishAuthRelationshipRequest(ctx context.Context, subjectType string, message events.AuthRelationshipRequest) error {
	var errors []error

	msg, err := e.conn.PublishAuthRelationshipRequest(ctx, subjectType, message)
	if err != nil {
		errors = append(errors, err)
	}

	if msg != nil {
		if msg.Error() != nil {
			errors = append(errors, err)
		}

		errors = append(errors, msg.Message().Errors...)
	}

	if len(errors) != 0 {
		e.logger.Error("error publishing request", "errors", errors)
	}

	return multierr.Combine(errors...)
}
