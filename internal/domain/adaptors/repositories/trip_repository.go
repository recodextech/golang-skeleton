package repositories

import "context"

type TripRepository interface {
	Exists(ctx context.Context, trip string) (exists bool, err error)
}
