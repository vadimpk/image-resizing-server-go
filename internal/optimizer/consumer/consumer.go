package consumer

import "context"

type Consumer interface {
	Consume(ctx context.Context, f func(body []byte, headers map[string]interface{})) error
	Close(ctx context.Context) chan struct{}
}
