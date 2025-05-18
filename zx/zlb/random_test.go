package zlb_test

import (
	"testing"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/zx/zlb"
)

var randomTargets = map[string]struct {
	found   bool
	name    string
	targets []*Target
}{
	"0 target": {
		found: false,
	},
	"1 target, inactive": {
		found: false,
		targets: []*Target{
			{name: "t0", id: 999000, weight: 1, active: false},
		},
	},
	"1 target, 0 weight": {
		found: false,
		targets: []*Target{
			{name: "t0", id: 999000, weight: 0, active: true},
		},
	},
	"1 target, non-0 weight": {
		found: true,
		name:  "t0",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 1, active: true},
		},
	},
	"2 targets, 0 weight": {
		found: false,
		targets: []*Target{
			{name: "t0", id: 999000, weight: 0, active: true},
			{name: "t1", id: 999111, weight: 0, active: true},
		},
	},
	"2 targets, non-0 weight": {
		found: true,
		name:  "t1",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 1, active: true},
			{name: "t1", id: 999111, weight: 2, active: true},
		},
	},
	"2 targets, equal weight": {
		found: true,
		name:  "t0",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 2, active: true},
			{name: "t1", id: 999111, weight: 2, active: true},
		},
	},
	"3 targets": {
		found: true,
		name:  "t2",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 3, active: true},
			{name: "t1", id: 999111, weight: 2, active: true},
			{name: "t2", id: 999222, weight: 1, active: true},
		},
	},
}

func TestRandom_Get(t *testing.T) {
	t.Parallel()
	for name, tc := range randomTargets {
		t.Run(name, func(t *testing.T) {
			lb := zlb.NewRandom(tc.targets...)
			lb.MaxRetry = 1
			tt, found := lb.Get(0)
			ztesting.AssertEqual(t, "found not match", tc.found, found)
			if !tc.found {
				return
			}
			ztesting.AssertEqual(t, "active status not match", true, tt.Active())
		})
	}
}

func TestRandomW_Get(t *testing.T) {
	t.Parallel()
	for name, tc := range randomTargets {
		t.Run(name, func(t *testing.T) {
			lb := zlb.NewRandomW(tc.targets...)
			tt, found := lb.Get(0)
			ztesting.AssertEqual(t, "found not match", tc.found, found)
			if !tc.found {
				return
			}
			ztesting.AssertEqual(t, "active status not match", true, tt.Active())
		})
	}
}
