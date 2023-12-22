package postgres

import (
	"database/sql"

	egos "github.com/finktek/egos/core"
)

type PostgresCheckpointStore struct {
	db *sql.DB
}

func NewPostgresCheckpointStore(db *sql.DB) egos.CheckpointStore {
	return &PostgresCheckpointStore{
		db: db,
	}
}

func (c *PostgresCheckpointStore) GetLastCheckpoint(checkpointID string) *egos.Checkpoint {
	var checkpoint egos.Checkpoint
	err := c.db.QueryRow("SELECT id, position FROM checkpoints WHERE id = $1", checkpointID).Scan(&checkpoint.ID, &checkpoint.Position)
	if err != nil {
		checkpoint.ID = checkpointID
		checkpoint.Position = 0
	}
	return &checkpoint
}

func (c *PostgresCheckpointStore) StoreCheckpoint(checkpoint *egos.Checkpoint) {
	_, err := c.db.Exec("INSERT INTO checkpoints (id, position) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET position = $2", checkpoint.ID, checkpoint.Position)
	if err != nil {
		panic(err)
	}
}
