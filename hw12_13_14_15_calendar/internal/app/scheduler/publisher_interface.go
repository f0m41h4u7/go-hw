package scheduler

type PublisherInterface interface {
	Connect() error
	Send([]byte) error
	Close() error
}
