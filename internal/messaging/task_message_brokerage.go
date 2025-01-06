package messaging

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
	"github.com/thomascastle/tarsk/internal/data"
)

type TaskMessageBrokerage struct {
	client *redis.Client
}

func NewTaskMessageBrokerage(client *redis.Client) *TaskMessageBrokerage {
	return &TaskMessageBrokerage{
		client: client,
	}
}

func (b *TaskMessageBrokerage) Created(ctx context.Context, task *data.Task) error {
	return b.publish(ctx, "tasks.event.created", task)
}

func (b *TaskMessageBrokerage) Deleted(ctx context.Context, id string) error {
	return b.publish(ctx, "tasks.event.deleted", id)
}

func (b *TaskMessageBrokerage) Updated(ctx context.Context, task *data.Task) error {
	return b.publish(ctx, "tasks.event.updated", task)
}

func (b *TaskMessageBrokerage) publish(ctx context.Context, channel string, event interface{}) error {
	var buf bytes.Buffer
	if e := json.NewEncoder(&buf).Encode(event); e != nil {
		return e
	}

	result := b.client.Publish(ctx, channel, buf.Bytes())
	if e := result.Err(); e != nil {
		return e
	}

	return nil
}
