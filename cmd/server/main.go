package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"test/internal/app"
	kafkaBroker "test/internal/broker/kafka"
	"test/internal/service"
	clickhouseStorage "test/internal/storage/clickhouse"
	postgresStorage "test/internal/storage/postgresql"
	redisStorage "test/internal/storage/redis"
	api "test/pkg/grpc"

	"google.golang.org/grpc"
)

var ctx = context.Background()

func main() {
	s := grpc.NewServer()

	// CLICKHOUSE
	ch, err := clickhouseStorage.New(ctx)
	if err != nil {
		fmt.Println("Main clickhouseStorage.New", err)
		return
	}
	defer func() {
		if err := ch.Close(); err != nil {
			fmt.Println("Main ch.Close", err)
		}
		fmt.Println("Clickhouse closed")
	}()
	if err := ch.Init(ctx); err != nil {
		fmt.Println("Main ch.Init", err)
		return
	}

	// POSTGRES
	postgres, err := postgresStorage.New(ctx)
	if err != nil {
		fmt.Println("Main postgresStorage.New", err)
		return
	}
	defer func() {
		if err := postgres.Close(); err != nil {
			fmt.Println("Main postgres.Close", err)
		}
		fmt.Println("Postgres closed")
	}()
	if err := postgres.Init(ctx); err != nil {
		fmt.Println("Main postgres.Init", err)
		return
	}

	// REDIS
	rdb, err := redisStorage.New()
	if err != nil {
		fmt.Println("Main redisStorage.New", err)
		return
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			fmt.Println("Main rdb.Close", err)
		}
		fmt.Println("Redis closed")
	}()

	// KAFKA
	kaf, err := kafkaBroker.New(ctx)
	if err != nil {
		fmt.Println("Main kafkaBroker.New", err)
		return
	}
	defer func() {
		if err := kaf.Close(); err != nil {
			fmt.Println("Main kaf.Close", err)
		}
		fmt.Println("Kafka closed")
	}()

	manager := service.New(postgres, rdb, kaf)

	srv := app.New(manager)

	api.RegisterMicroserviceServer(s, srv)

	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Server started")

	if err := s.Serve(l); err != nil {
		log.Fatal(err)
	}
}
