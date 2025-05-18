package zrate

import (
	"context"
	"sync"
)

var (
	_ Token   = NoopToken(true)
	_ Token   = &token{}
	_ Limiter = &ConcurrentLimiter{}
	_ Limiter = &BucketLimiter{}
	_ Limiter = &LeakyBucketLimiter{}
)

const (
	// TokenOK always reports true with [Token.OK] and nil error.
	TokenOK = NoopToken(true)
	// TokenNG always reports false with [Token.OK] and nil error.
	TokenNG = NoopToken(false)
)

// Limiter provides rate limiting mechanism.
// See also [golang.org/x/time/rate].
type Limiter interface {
	AllowNow() Token
	WaitNow(context.Context) Token
}

// Token represents limiter tokens.
type Token interface {
	// OK returns if the token is valid or not.
	// If false, callers' must not proceed their process.
	OK() bool
	// Release releases the token.
	// It's depends on the implementation, or algorithms, that calling
	// Release is mandatory or not.
	// So, it's recommended for callers' to call Release() for any implementation.
	// It's safe to call Release multiple times.
	Release()
	// Err returns error occurred while obtaining tokens.
	Err() error
}

// NoopLimiter always reports the fixed value of [TokenOK] or [TokenNG].
// NoopLimiter implements [Limiter] interface
type NoopLimiter bool

func (lim NoopLimiter) AllowNow() Token {
	if lim {
		return TokenOK
	}
	return TokenNG
}

func (lim NoopLimiter) WaitNow(_ context.Context) Token {
	if lim {
		return TokenOK
	}
	return TokenNG
}

// NoopToken always reports the fixed value of true or false.
// Release() does nothing and Err() always returns nil.
// NoopToken implements [Token] interface
type NoopToken bool

func (t NoopToken) OK() bool {
	return bool(t)
}

func (t NoopToken) Release() {
	// Nothing to do.
}

func (t NoopToken) Err() error {
	return nil
}

// token is the default token that implements [Token].
type token struct {
	once        sync.Once
	releaseFunc func()
	ok          bool
	err         error
}

func (t *token) OK() bool {
	return t.ok
}

func (t *token) Release() {
	if t.releaseFunc != nil {
		t.once.Do(t.releaseFunc)
	}
}

func (t *token) Err() error {
	return t.err
}
