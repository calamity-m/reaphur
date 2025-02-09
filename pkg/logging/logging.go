package logging

import (
	"context"
	"io"
	"log/slog"

	"github.com/calamity-m/reaphur/pkg/bindings"
)

type CustomHandler struct {
	slog.Handler
	recordRequestId bool
}

type CustomHandlerCfg struct {
	// Level to set for the logger, generic shit
	Level slog.Level

	// Whether to print logs in a structured json format
	// or in pretty text
	Structed bool

	// Whether to display source of the log line
	AddSource bool

	// Whether to grab the request id from the request's context
	// variable and append it to log lines created by this
	// handler
	RecordRequestId bool

	// Static attributes that will be appended to every
	// log line created by this handler
	StaticAttributes []slog.Attr
}

func (c *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	if c.recordRequestId {
		reqId, ok := ctx.Value(bindings.RequestIDKey{}).(string)

		if ok {
			r.AddAttrs(slog.String("request-id", string(reqId)))
		} else {
			r.AddAttrs(slog.String("request-id", "unknown"))
		}
	}

	// Stop intercepting and continue on
	err := c.Handler.Handle(ctx, r)
	return err
}

func NewCustomizedHandler(w io.Writer, cfg *CustomHandlerCfg) *CustomHandler {
	handler := &CustomHandler{}

	if cfg == nil {
		cfg = &CustomHandlerCfg{}
	}

	handler.recordRequestId = cfg.RecordRequestId

	handlerOpts := &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: cfg.AddSource,
	}

	if cfg.Structed {
		handler.Handler = slog.NewJSONHandler(w, handlerOpts)
	} else {
		handler.Handler = slog.NewTextHandler(w, handlerOpts)
	}

	handler.Handler = handler.Handler.WithAttrs(cfg.StaticAttributes)

	return handler
}
