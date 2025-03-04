
```
go-machine-events
├─ Dockerfile
├─ README.md
├─ cmd
│  ├─ app
│  │  └─ main.go
│  ├─ consumer
│  │  └─ main.go
│  └─ producer
│     └─ main.go
├─ config
│  ├─ config.go
│  └─ config.yaml
├─ docker-compose.yml
├─ go.mod
├─ go.sum
├─ internal
│  ├─ events
│  │  ├─ refill_event.go
│  │  └─ sale_event.go
│  ├─ machine
│  │  ├─ machine.go
│  │  └─ repository.go
│  ├─ models
│  ├─ pubsub
│  │  ├─ publisher.go
│  │  ├─ pubsub.go
│  │  └─ subscriber.go
│  └─ services
│     ├─ refill_service.go
│     └─ sale_service.go
├─ pkg
│  ├─ kafka
│  │  ├─ consumer.go
│  │  ├─ kafka_client.go
│  │  └─ producer.go
│  └─ logger
│     └─ logger.go
└─ test
   ├─ machine_test.go
   └─ pubsub_test.go

```