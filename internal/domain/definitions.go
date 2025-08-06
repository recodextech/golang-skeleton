package domain

const (
	TripCreateHandler = `http.handler.trip-create`
)

type ContextKey string

const (
	DateLayout     = "2006-01-02"
	DateTimeLayout = "2006-01-02 15:04:05"
)

const (
	ContextKeyAccountID ContextKey = `account-id`
	ContextKeyTraceID   ContextKey = `trace-id`
	ContextKeyUserID    ContextKey = `user-id`
	ContextKeyUserType  ContextKey = `triggered-user-type`
	ContextKeyStreamID  ContextKey = "stream-id"
	ContextKeyTimeZone  ContextKey = "time-zone"
)

func (c ContextKey) String() string {
	return string(c)
}
