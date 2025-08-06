package events

import (
	"context"
	"golang-skeleton/internal/domain"
	"golang-skeleton/pkg/uuid"
	"time"
)

type Meta struct {
	Type      string    `json:"type"`
	EventID   string    `json:"event_id"`
	Expiry    int       `json:"expiry"`
	CreatedAt int64     `json:"created_at"`
	Version   int64     `json:"version"`
	TraceID   string    `json:"trace_id"`
	CreatedBy string    `json:"created_by"`
	AccountID uuid.UUID `json:"account_id"`
}

type Event interface {
	GetMeta() Meta
}

type Header string

func NewMeta() Meta {
	return Meta{
		EventID:   uuid.New().String(),
		Version:   1,
		CreatedAt: time.Now().UTC().UnixNano() / 1e6,
		TraceID:   uuid.New().String(),
	}
}

func NewMetaContext(ctx context.Context) Meta {
	meta := NewMeta()
	return MetaUpdate(ctx, meta)
}

func MetaUpdate(ctx context.Context, meta Meta) Meta {
	meta.TraceID = extractFromContext(ctx, domain.ContextKeyTraceID.String(), uuid.New().String())
	meta.CreatedBy = extractFromContext(ctx, domain.ContextKeyUserID.String(), "")
	meta.AccountID = ctx.Value(domain.ContextKeyAccountID.String()).(uuid.UUID)

	return meta
}

func extractFromContext(ctx context.Context, key interface{}, defaultVal string) string {
	val, ok := ctx.Value(key).(string)
	if !ok {
		val = defaultVal
	}

	return val
}

const (
	EventHeaderTraceID   Header = `trace_id`
	EventHeaderAccountID Header = `account_id`
)
