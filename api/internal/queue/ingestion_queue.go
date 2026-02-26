package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type IngestionQ struct {
	rdb *redis.Client
}

type IngestionJob struct {
	ID       int32  `json:"id"`
	RefID    int32  `json:"ref_id"`
}

const (
	queueName = "ingestion_queue"
)

func NewIngestionQ(rdb *redis.Client, maxLen int) (*IngestionQ, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := rdb.XTrimMaxLen(ctx, queueName, int64(maxLen)).Err()
	if err != nil {
		return nil, err
	}

	return &IngestionQ{
		rdb: rdb,
	}, nil
}

func (q *IngestionQ) Enqueue(ctx context.Context, j *IngestionJob) error {
	bytes, err := json.Marshal(j)
	if err != nil {
		return err
	}

	err = q.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: queueName,
		ID:     "*",
		Values: [][]byte{[]byte("data"), bytes},
	}).Err()
	if err != nil {
		return err
	}
	return nil
}
