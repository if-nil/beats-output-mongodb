package mongoout

import (
	"context"
	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/logp"
	"github.com/elastic/beats/v7/libbeat/outputs"
	"github.com/elastic/beats/v7/libbeat/publisher"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type client struct {
	collection     *mongo.Collection
	log            *logp.Logger
	host           string
	dbName         string
	collectionName string
	client         *mongo.Client
	observer       outputs.Observer
	info           beat.Info
	timeout        time.Duration
}

var _ outputs.Client = &client{}

func newClient(host, dbName, collectionName string, observer outputs.Observer, info beat.Info, timeout time.Duration) (*client, error) {
	c := &client{
		log:            logp.NewLogger("mongodb"),
		host:           host,
		dbName:         dbName,
		collectionName: collectionName,
		observer:       observer,
		info:           info,
		timeout:        timeout,
	}
	if err := c.Connect(); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *client) Close() error {
	return c.client.Disconnect(context.TODO())
}

func (c *client) Publish(ctx context.Context, batch publisher.Batch) error {
	if batch == nil {
		panic("no batch")
	}

	events := batch.Events()
	c.observer.NewBatch(len(events))

	docs, _ := serializeEvents(c, events)
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()
	_, err := c.collection.InsertMany(ctx, docs)
	if err != nil {
		c.log.Errorf("Failed to insert events into MongoDB: %v", err)
		c.observer.Failed(len(events))
		batch.RetryEvents(events)
		return err
	}
	batch.ACK()
	return nil
}

func (c *client) String() string {
	return "mongodb"
}

func (c *client) Connect() error {
	c.log.Debug("Connecting to MongoDB")
	var err error
	c.client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(c.host))
	if err != nil {
		return err
	}
	c.collection = c.client.Database(c.dbName).Collection(c.collectionName)
	return nil
}

func serializeEvents(c *client, events []publisher.Event) ([]interface{}, error) {
	var docs []interface{}
	for i := range events {
		e := events[len(events)-1-i]
		var m = bson.M(e.Content.Fields)
		m["_timestamp"] = e.Content.Timestamp
		m["_timezone"] = e.Content.Timestamp.Location().String()
		m["_metadata"] = bson.M{
			"beat":    c.info.Beat,
			"uuid":    c.info.ID.String(),
			"version": c.info.Version,
		}
		docs = append(docs, m)
	}
	return docs, nil
}
