# accounts gocqrs example

Start PostgreSQL, Zookeeper, and Kafka services.
```shell script
$ docker-compose up
```

In another terminal, create the 'Account' topic where domain events will be published.
```shell script
$ docker run --rm --network accounts_default confluentinc/cp-kafka:6.0.0 kafka-topics --bootstrap-server kafka:29092 --create --topic Account --partitions 5 --replication-factor 1
```

Set environment variables for connecting to Postgres and Kafka.
```shell script
$ export DATABASE_URL="postgresql://localhost:5432/accounts?user=accounts&password=accounts"
$ export KAFKA_BOOTSTRAP_SERVERS="localhost"
```

Build and run the app.
```shell script
$ go build
$ ./accounts
```

In a third terminal, run a kafka console consumer to see the domain events as they're being written. Also use this to
verify that events for the same "entityId" are published on the same partition to guarantee ordering.
```shell script
$ docker run -it --network=host edenhill/kafkacat:1.6.0 -b localhost -C -t Account -f '%t %p @ %o: %s\n'
```

In a fourth terminal, run the following curl commands to generate some events.

Create an account.
```shell script
$ curl --location --request POST 'localhost:8080/accounts/' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "name": "Jane Doe",
      "startingBalance": 120
  }'

{"id":"bului7u602fmfhogc05g","version":1,"name":"Jane Doe","balance":120,"status":"Open"}
```

Fetch the bewly created account.
```shell script
$ curl --location --request GET 'localhost:8080/accounts/bului7u602fmfhogc05g'
```

Credit the account.
```shell script
$ curl --location --request PUT 'localhost:8080/accounts/bului7u602fmfhogc05g' \
  --header 'Content-Type: application/json' \
  --data-raw '{
      "creditAmount": 10
  }'
```

Note that all events for the same entity are published on the same Kafka partition and the entity version is updated
every time to allow for optimistic locking.
