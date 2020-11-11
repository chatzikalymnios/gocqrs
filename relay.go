package gocqrs

type Relay interface {
	Relay()
	Publisher() Publisher
}
