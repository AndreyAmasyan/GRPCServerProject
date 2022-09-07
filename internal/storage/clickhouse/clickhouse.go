package clickhouseStorage

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mailru/go-clickhouse/v2"
)

// ClickhouseStorage ...
type ClickhouseStorage struct {
	db *sql.DB
}

func New(ctx context.Context) (*ClickhouseStorage, error) {

	connect, err := sql.Open(
		"chhttp",
		"http://"+os.Getenv("CH_HOST")+":"+os.Getenv("CH_PORT")+"/"+os.Getenv("CH_DB"),
	)
	if err != nil {
		fmt.Println("ClickhouseStorage New sql.Open", err)
		return nil, err
	}

	if err := connect.PingContext(ctx); err != nil {
		fmt.Println("ClickhouseStorage New db.PingContext", err)
		return nil, err
	}

	return &ClickhouseStorage{
		db: connect,
	}, nil
}

func (s *ClickhouseStorage) Init(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS data (
			id UInt32,
			name String,
			action_time DateTime
		) ENGINE = MergeTree()
			ORDER BY id;
	`)
	if err != nil {
		fmt.Println("Can't create table data in clickhouse", err)
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS data_kafka(
			id UInt32,
			name String,
			action_time DateTime
		) ENGINE = Kafka
			SETTINGS kafka_broker_list = ?,
		  	kafka_topic_list = 'data_topic',
			kafka_group_name = 'data_group_name',
		  	kafka_format = 'CSV';
	`, os.Getenv("KAFKA_HOST_PORT"))

	if err != nil {
		fmt.Println("Can't create table data_kafka in clickhouse", err)
		return err
	}

	_, err = s.db.ExecContext(ctx, `
		CREATE MATERIALIZED VIEW IF NOT EXISTS data_mv TO data AS
		SELECT
		  	id,
		  	name,
		  	action_time
		FROM data_kafka;
	`)

	if err != nil {
		fmt.Println("Can't create materialized view data_mv in clickhouse", err)
	}

	return nil
}

func (s *ClickhouseStorage) Close() error {
	if err := s.db.Close(); err != nil {
		fmt.Println("ClickhouseStorage Close s.db.Close", err)
		return err
	}

	return nil
}
