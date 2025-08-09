package writers

import (
	"context"
	"encoding/json"
	"golang-skeleton/internal/domain/adaptors"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/domain/events"
	"golang-skeleton/internal/http/responses"
	"golang-skeleton/pkg/errors"
	"net/http"

	"github.com/recodextech/container"

	"github.com/tryfix/krouter"
)

type TripResponseWriter struct {
	log adaptors.Logger
}

func (t *TripResponseWriter) Response(_ context.Context, w http.ResponseWriter, _ *http.Request,
	payload krouter.HttpPayload) error {
	var err error
	out := payload.Body.(events.TripCreate)
	response := responses.SuccessResponse{}
	response.Data.ID = out.Payload.PassengerID

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resBody, err := json.Marshal(response)
	if err != nil {
		return ResponseWriterError{errors.Wrap(err, errorWritingResponse)}
	}
	_, err = w.Write(resBody)
	if err != nil {
		return ResponseWriterError{errors.Wrap(err, errorWritingResponse)}
	}

	return nil
}

func (t *TripResponseWriter) Init(container container.Container) error {
	t.log = container.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed(
		`responses.trip`))
	return nil
}
