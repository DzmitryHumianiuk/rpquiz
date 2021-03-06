package bot

import (
	"context"
	"github.com/pkg/errors"
)

// ErrUnknownIntent general error to be thrown in case intent not found
var ErrUnknownIntent = errors.New("intent is unknown")

const acceptableConfidence = 0.5

//IntentNameDispatcher is a composite Handler for DialogFlow. Dispatches incoming request over specific intent handlers
func IntentNameDispatcher(handlers map[string]Handler, fallback Handler) Handler {
	return HandlerFunc(func(ctx context.Context, rq Request) ([]*Response, error) {
		var handler Handler

		if irq, ok := rq.(*IntentRequest); ok {
			if irq.Confidence < acceptableConfidence {
				//intent isn't recognized. No reason to search for intent handler at all
				handler = fallback
			}
			if h, ok := handlers[irq.Intent]; ok {
				//intent is recognized and handler is implemented
				handler = h
			} else {
				//intent is recognized but handler not implemented
				handler = fallback
			}

		} else {
			//intent is recognized but handler not implemented
			handler = fallback
		}
		return handler.Handle(ctx, rq)
	})
}
