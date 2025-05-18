package zuid

import (
	"context"
)

type ctxKey struct{ string } // Context key for saving unique IDs.

// ContextWithID save an unique ID in the context with given key.
// A new context created with [context.Background] is used when
// nil context was given. Note that calling ContextWithID multiple times
// overwrites the existing unique ID.
// Use [FromContext] to extract stored id from the context.
func ContextWithID(ctx context.Context, key, id string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ctxKey{key}, id)
}

// FromContext returns an unique ID extracted from the context.
// Empty string will be returned if no unique ID found
// in the context or the context is nil.
// Use [ContextWithID] to save an id to the context.
func FromContext(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}
	id, _ := ctx.Value(ctxKey{key}).(string)
	return id
}
