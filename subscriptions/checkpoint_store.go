package subscriptions

type CheckpointStore interface {
	GetLastCheckpoint(subscriptionId string) Checkpoint
	StoreCheckpoint(checkpoint Checkpoint)
}


type Checkpoint struct {
	ID 			string
	Position 	uint64
}