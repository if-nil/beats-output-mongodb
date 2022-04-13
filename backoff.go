package mongoout

import (
	"context"
	b "github.com/elastic/beats/v7/libbeat/common/backoff"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/publisher"
	"time"
)

type backoffClient struct {
	client  *client
	done    chan struct{}
	backoff b.Backoff
}

var _ outputs.NetworkClient = &backoffClient{}

func newBackoffClient(client *client, init, max time.Duration) *backoffClient {
	done := make(chan struct{})
	backoff := b.NewEqualJitterBackoff(done, init, max)
	return &backoffClient{
		client:  client,
		done:    done,
		backoff: backoff,
	}
}

func (b *backoffClient) Close() error {
	close(b.done)
	return b.client.Close()
}

func (b *backoffClient) Publish(ctx context.Context, batch publisher.Batch) error {
	err := b.client.Publish(ctx, batch)
	if err != nil {
		b.backoff.Wait()
	}
	return err
}

func (b *backoffClient) String() string {
	return b.client.String()
}

func (b *backoffClient) Connect() error {
	return nil
}
