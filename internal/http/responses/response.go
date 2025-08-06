package responses

import (
	"context"
	"net/http"

	"github.com/tryfix/krouter"
)

type GenerateResponse interface {
	Response(ctx context.Context, w http.ResponseWriter, r *http.Request, payload krouter.HttpPayload) error
}
