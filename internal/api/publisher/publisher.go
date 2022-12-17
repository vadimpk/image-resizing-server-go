package publisher

import "context"

type Publisher interface {
	Publish(ctx context.Context, body []byte, headers map[string]interface{}) error
	Close(ctx context.Context) chan struct{}
}
