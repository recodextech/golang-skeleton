package handlers

import (
	"context"
	"encoding/json"
	"golang-skeleton/internal/app"
	"golang-skeleton/internal/domain/adaptors"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/http/request"
	"golang-skeleton/pkg/container"
	"net/http"
)

// ErrorHandler contains all error handling, formatting and presenting functionality for the http layer.
type ErrorHandler struct {
	log adaptors.Logger
}

// Init initializes a new instance of the error handler.
func (hlr *ErrorHandler) Init(c container.Container) error {
	hlr.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(
		adaptors.LoggerPrefixed("ErrorHandler"))

	return nil
}

// Handle handles all errors globally.
func (hlr *ErrorHandler) Handle(ctx context.Context, w http.ResponseWriter, e error) error {
	// log the error
	hlr.log.ErrorContext(ctx, e.Error())
	traceID := ctx.Value(request.TraceID.String()).(string)

	resError := format(e)
	resError.Trace = traceID

	if app.DebugMode() {
		resError.Debug = e.Error()
	}

	payload, err := json.Marshal(resError)
	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(resError.Status)

	_, err = w.Write(payload)
	if err != nil {
		return err
	}

	return nil
}
