package subscriptions

import "github.com/finktek/eventum/subscriptions"

type MongoCheckpointStore struct {
}

func NewMongoDbCheckpointStore() (*MongoCheckpointStore, error) {
	return &MongoCheckpointStore{}, nil
}

func (cs MongoCheckpointStore) GetLastCheckpoint(subscriptionId string) subscriptions.Checkpoint {
	return subscriptions.Checkpoint{ID: "csID", Position: 439661}
}

func (cs MongoCheckpointStore) StoreCheckpoint(checkpoint subscriptions.Checkpoint) {

}