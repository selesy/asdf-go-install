// Package logger provides a caching slog.Handler for use before the
// logging system is fully configured.
package logger

import (
	"context"
	"log/slog"
	"strings"
)

var _ slog.Handler = (*CachingHandler)(nil)

// CachingHandler collects log records that have to be generated before
// an application has configured a logger that creates output.
//
// This should be an extremely short life-time so that log messages are
// not lost if the application crashes.
type CachingHandler struct {
	recs  *[]slog.Record
	attrs []slog.Attr
	grps  []string
}

// NewCachingHandler creates a Caching slog.Handler.
func NewCachingHandler() *CachingHandler {
	return &CachingHandler{
		recs: &[]slog.Record{},
	}
}

// Enabled implements slog.Handler.
func (h *CachingHandler) Enabled(ctx context.Context, lvl slog.Level) bool {
	return true
}

// Handle implements slog.Handler.
func (h *CachingHandler) Handle(ctx context.Context, rec slog.Record) error {
	// Create a clone of the record without its attributes
	clone := slog.Record{
		Time:    rec.Time,
		Message: rec.Message,
		PC:      rec.PC,
		Level:   rec.Level,
	}

	// Add the parent handlers attributes to the record.
	clone.AddAttrs(h.attrs...)

	// Prepend the group(s), if any, to each of the attributes in the
	// incoming record and add the updated attribute to the cloned
	// record.
	rec.Attrs(func(a slog.Attr) bool {
		a.Key = strings.Join(append(h.grps, a.Key), ".")
		clone.AddAttrs(a)

		return true
	})

	// Append the new record to the cache.
	*h.recs = append(*h.recs, clone)

	return nil
}

// WithAttrs implements slog.Handler.
func (h *CachingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// The existing handler's groups need to be prepended to the incoming
	// attributes' keys.
	for i, attr := range attrs {
		attrs[i].Key = strings.Join(append(h.grps, attr.Key), ".")
	}

	// We have to provide a pointer to the parent Handler's records, but
	// an updated copy of the parent's attributes and an unchanged copy
	// of the parent's groups.
	return &CachingHandler{
		recs:  h.recs,
		attrs: append(h.attrs, attrs...),
		grps:  h.grps,
	}
}

// WithGroup implements slog.Handler.
func (h *CachingHandler) WithGroup(grp string) slog.Handler {
	// We have to provide a pointer to the parent Handler's records, but
	// an unchanged copy of the paren's attributes and an updated copy
	// of the parent's groups.
	return &CachingHandler{
		recs:  h.recs,
		attrs: h.attrs,
		grps:  append(h.grps, grp),
	}
}

// Transfer copies the cached log records to the "final" slog.Logger.
func (h *CachingHandler) Transfer(ctx context.Context, log *slog.Logger) error {
	for _, rec := range *h.recs {
		if err := log.Handler().Handle(ctx, rec); err != nil {
			return err
		}
	}

	return nil
}
