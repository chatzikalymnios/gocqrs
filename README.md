# gocqrs

Small side project to get familiar with `golang`, event sourcing, cqrs and the
[transactional outbox pattern](https://microservices.io/patterns/data/transactional-outbox.html) for publishing domain
events correctly.

This project includes a small event store built on top of PostgreSQL, although the design allows for a different storage
solution by creating new concrete implementations for a few interfaces (`AggregateStore`, `Relay`). It also uses Kafka to
publish the domain events.

A small example application can be seen in [`example/accounts`](example/accounts/README.md).
