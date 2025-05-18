package zuid_test

import (
	"context"
	"testing"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/zx/zuid"
)

func TestContextWithID(t *testing.T) {
	t.Parallel()
	t.Run("nil context", func(t *testing.T) {
		ctx := zuid.ContextWithID(nil, "key", "test-uid")
		uid := zuid.FromContext(ctx, "key")
		ztesting.AssertEqual(t, "uid not match", "test-uid", uid)
	})
	t.Run("non-nil context", func(t *testing.T) {
		ctx := zuid.ContextWithID(context.Background(), "key", "test-uid")
		uid := zuid.FromContext(ctx, "key")
		ztesting.AssertEqual(t, "uid not match", "test-uid", uid)
	})
	t.Run("key not match", func(t *testing.T) {
		ctx := zuid.ContextWithID(context.Background(), "key1", "test-uid")
		uid := zuid.FromContext(ctx, "key2")
		ztesting.AssertEqual(t, "uid not match", "", uid)
	})
}

func TestFromContext(t *testing.T) {
	t.Parallel()
	t.Run("nil context", func(t *testing.T) {
		uid := zuid.FromContext(nil, "key")
		ztesting.AssertEqual(t, "uid not match", "", uid)
	})
	t.Run("no uid", func(t *testing.T) {
		uid := zuid.FromContext(context.Background(), "key")
		ztesting.AssertEqual(t, "uid not match", "", uid)
	})
	t.Run("non-nil context", func(t *testing.T) {
		ctx := zuid.ContextWithID(context.Background(), "key", "test-uid")
		uid := zuid.FromContext(ctx, "key")
		ztesting.AssertEqual(t, "uid not match", "test-uid", uid)
	})
	t.Run("key not match", func(t *testing.T) {
		ctx := zuid.ContextWithID(context.Background(), "key1", "test-uid")
		uid := zuid.FromContext(ctx, "key2")
		ztesting.AssertEqual(t, "uid not match", "", uid)
	})
}
