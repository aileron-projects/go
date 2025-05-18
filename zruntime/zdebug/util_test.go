package zdebug_test

import (
	"testing"

	"github.com/aileron-projects/go/zruntime/zdebug"
)

func TestGoDebugEnv(t *testing.T) {
	// t.Parallel() // This test cannot be parallel.

	envPrefix := "zruntime-zdebug-TestGoDebugEnv"
	t.Setenv("GODEBUG", envPrefix+"-01=test-value")

	t.Run("key not found", func(t *testing.T) {
		key := envPrefix + "-00"
		v, ok := zdebug.GoDebugEnv(key)
		if ok {
			t.Errorf("unexpected env found. key=%s. got:%s", key, v)
		}
		if v != "" {
			t.Errorf("unexpected env found. key=%s. got:%s", key, v)
		}
	})

	t.Run("key found", func(t *testing.T) {
		key := envPrefix + "-01"
		v, ok := zdebug.GoDebugEnv(key)
		if !ok {
			t.Errorf("env not found. key=%s. want:test-value got:%s", key, v)
		}
		if v != "test-value" {
			t.Errorf("unexpected env value. key=%s. want:test-value got:%s", key, v)
		}
	})
}
