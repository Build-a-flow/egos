# Egos - Event Sourcing library for GO

## Table of Contents
- [Checkpoint Migration](#checkpoint-migration)

## Checkpoint Migration
To perform the checkpoint migration, execute the following SQL query:

### For Postgres:
```DROP TABLE IF EXISTS "checkpoints";
CREATE TABLE "checkpoints" (
    "id" text NOT NULL,
    "position" bigint NOT NULL,
    CONSTRAINT "checkpoints_id" UNIQUE ("id")
);```

