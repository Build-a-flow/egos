package subscriptions

import subscriptions "github.com/build-a-flow/egos/subscriptions"

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
