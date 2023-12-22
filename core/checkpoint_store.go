package egos

type CheckpointStore interface {
	GetLastCheckpoint(checkpointID string) *Checkpoint
	StoreCheckpoint(checkpoint *Checkpoint)
}

type Checkpoint struct {
	ID       string
	Position uint64
}
