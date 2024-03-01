package sesm

import (
	"context"
	"time"
)

type Store interface {
	StoreFind(context.Context, string) (int, time.Time, error)
	StoreCommit(context.Context, string, int, time.Time) error
	StoreDeleteAll(context.Context, int) error
	StoreDelete(context.Context, string) error
}
