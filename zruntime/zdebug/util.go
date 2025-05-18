package zdebug

import (
	"os"
	"strings"
)

// GoDebugEnv returns the value of GODEBUG
// which has the given key. When GODEBUG=foo=bar,alice=bob
// GoDebugEnv("foo") returns "bar", GoDebugEnv("alice") returns "bob".
// GoDebugEnv returns true when the given key was found in the GODEBUG env
// and returns false when not found.
// GoDebugEnv does not trim or remove any spaces.
// So, if GODEBUG="foo=bar, alice=bob" then GoDebugEnv("alice") returns false.
func GoDebugEnv(key string) (string, bool) {
	debugs := strings.Split(os.Getenv("GODEBUG"), ",")
	key = key + "="
	for _, debug := range debugs {
		if strings.HasPrefix(debug, key) {
			return strings.TrimPrefix(debug, key), true
		}
	}
	return "", false
}
