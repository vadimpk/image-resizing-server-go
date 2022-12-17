package consumer

type QueueConfig struct {
	Name          string
	Durable       bool
	DeleteUnused  bool
	Exclusive     bool
	NoWait        bool
	PrefetchCount int
}
