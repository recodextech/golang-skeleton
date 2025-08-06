package http

import (
	"context"
	"encoding/json"
	"golang-skeleton/internal/domain/adaptors"
	"golang-skeleton/internal/domain/application"
	"golang-skeleton/internal/http/handlers"
	"golang-skeleton/internal/http/request"
	"golang-skeleton/internal/http/responses"
	"golang-skeleton/pkg/container"
	"golang-skeleton/pkg/uuid"
	"net/http"
	"strconv"

	"github.com/tryfix/krouter"
	"github.com/tryfix/log"
	"github.com/tryfix/metrics"

	"github.com/gorilla/mux"
)

const (
	base10    = 10
	bitSize64 = 64
)

type Router struct {
	router       *mux.Router
	log          adaptors.Logger
	krouter      *krouter.Router
	routerConfig *KRouterConf
}

func (r *Router) Init(c container.Container) error {
	r.log = c.Resolve(application.ModuleLogger).(adaptors.Logger).NewLog(adaptors.LoggerPrefixed("router"))
	r.routerConfig = c.GetGlobalConfig(application.ModuleHTTPRouter).(*KRouterConf)
	errorHandle := c.Resolve(handlers.ModuleErrorHandler).(*handlers.ErrorHandler)

	kafRouter, err := krouter.NewRouter(krouter.Config{
		BootstrapServers: r.routerConfig.BootstrapServers,
		RouterTopic:      r.routerConfig.Topic,
		ConsumerGroup:    r.routerConfig.ApplicationID,
	},
		krouter.WithLogger(c.Resolve(application.ModuleBaseLogger).(log.Logger).NewLog(log.Prefixed(`async-router`))),
		krouter.WithParamType(request.ParamTypeAppUUID, func(v string) (interface{}, error) {
			return uuid.Parse(v)
		}),
		krouter.WithParamType(request.ParamTypeAppInt64, func(v string) (interface{}, error) {
			return strconv.ParseInt(v, base10, bitSize64)
		}),
		krouter.WithParamType(request.ParamTypeAppString, func(v string) (interface{}, error) {
			return v, nil
		}),
		krouter.WithContextExtractor(func(req *http.Request) context.Context {
			userID, err := uuid.Parse(req.Header.Get(request.HeaderUserID.String()))
			if err != nil {
				userID = uuid.Nil
			}
			ctx := context.WithValue(req.Context(), request.UserID.String(), userID)

			accID, err := uuid.Parse(req.Header.Get(request.HeaderAccountID.String()))
			if err != nil {
				accID = uuid.Nil
			}
			ctx = context.WithValue(ctx, request.AccountID.String(), accID)

			traceID, err := uuid.Parse(req.Header.Get(request.HeaderTraceID.String()))
			if err != nil {
				traceID = uuid.Nil
			}

			ctx = context.WithValue(ctx, request.TraceID.String(), traceID.String())
			// extract custom headers set from the api-gateway
			ctx = r.applyXHeaders(ctx, req)

			return ctx
		}),
		krouter.WithErrorHandlerFunc(func(ctx context.Context, w http.ResponseWriter, r *http.Request, err error) error {
			return errorHandle.Handle(ctx, w, err)
		}),
		krouter.WithMetricsReporter(c.Resolve(application.ModuleStreamReporter).(metrics.Reporter)),
	)
	if err != nil {
		return err
	}

	r.router = mux.NewRouter()
	r.krouter = kafRouter

	// routes
	initRoutes(r, c)

	return nil
}

func (r *Router) Run() error {
	return r.krouter.Start()
}
func (r *Router) Stop() error {
	return nil
}

// applyXHeaders to retrieve particular set of headers from the api gateway layer
func (r Router) applyXHeaders(ctx context.Context, req *http.Request) context.Context {
	m := make(map[string]string)
	headerArr := r.routerConfig.XHeaders

	for _, k := range headerArr {
		m[k] = req.Header.Get(k)
	}

	return context.WithValue(ctx, request.XHeaders, m)
}

func successResponse(mapper func(payload krouter.HttpPayload) string, httpStatus int) krouter.SuccessHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, payload krouter.HttpPayload) error {
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(httpStatus)
		if httpStatus == http.StatusNoContent {
			return nil
		}

		res := responses.TripCreateResponse{}

		resBody, err := json.Marshal(res)
		if err != nil {
			return err
		}
		_, err = w.Write(resBody)
		return err
	}
}

func traceIDGenerate() string {
	return uuid.New().String()
}
