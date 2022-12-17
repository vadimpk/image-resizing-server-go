package publisher

import "context"

type Queue interface {
	PublishImage(ctx context.Context, body []byte, headers map[string]interface{})
}
