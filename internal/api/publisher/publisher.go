package publisher

import "context"

type Publisher interface {
	PublishImage(ctx context.Context, body []byte, headers map[string]interface{})
}
