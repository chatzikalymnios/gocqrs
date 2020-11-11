package gocqrs

type Publisher interface {
	Publish()
	Events() chan *Event
	Err() chan error
}
