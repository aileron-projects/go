package zlb_test

import (
	"testing"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/zx/zlb"
)

var priTargets = map[string]struct {
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
			{name: "t0", id: 0, weight: 1, active: false},
		},
	},
	"1 target, 0 weight": {
		found: false,
		targets: []*Target{
			{name: "t0", id: 0, weight: 0, active: true},
		},
	},
	"1 target, non-0 weight": {
		found: true,
		name:  "t0",
		targets: []*Target{
			{name: "t0", id: 0, weight: 1, active: true},
		},
	},
	"2 targets, 0 weight": {
		found: false,
		targets: []*Target{
			{name: "t0", id: 0, weight: 0, active: true},
			{name: "t1", id: 1, weight: 0, active: true},
		},
	},
	"2 targets, non-0 weight": {
		found: true,
		name:  "t1",
		targets: []*Target{
			{name: "t0", id: 0, weight: 1, active: true},
			{name: "t1", id: 1, weight: 2, active: true},
		},
	},
	"2 targets, equal weight": {
		found: true,
		name:  "t0",
		targets: []*Target{
			{name: "t0", id: 0, weight: 2, active: true},
			{name: "t1", id: 1, weight: 2, active: true},
		},
	},
	"3 targets, weights 321": {
		found: true,
		name:  "t0",
		targets: []*Target{
			{name: "t0", id: 0, weight: 3, active: true},
			{name: "t1", id: 1, weight: 2, active: true},
			{name: "t2", id: 2, weight: 1, active: true},
		},
	},
	"3 targets, weights 231": {
		found: true,
		name:  "t1",
		targets: []*Target{
			{name: "t0", id: 0, weight: 2, active: true},
			{name: "t1", id: 1, weight: 3, active: true},
			{name: "t2", id: 2, weight: 1, active: true},
		},
	},
	"3 targets, weights 123": {
		found: true,
		name:  "t2",
		targets: []*Target{
			{name: "t0", id: 0, weight: 1, active: true},
			{name: "t1", id: 1, weight: 2, active: true},
			{name: "t2", id: 2, weight: 3, active: true},
		},
	},
	"3 targets, contains inactive": {
		found: true,
		name:  "t1",
		targets: []*Target{
			{name: "t0", id: 0, weight: 1, active: true},
			{name: "t1", id: 1, weight: 2, active: true},
			{name: "t2", id: 2, weight: 3, active: false},
		},
	},
}

func TestPriority_Get(t *testing.T) {
	t.Parallel()
	for name, tc := range priTargets {
		t.Run(name, func(t *testing.T) {
			lb := zlb.NewRoundRobin(tc.targets...)
			tt, found := lb.Get(0)
			ztesting.AssertEqual(t, "found not match", tc.found, found)
			if !tc.found {
				return
			}
			ztesting.AssertEqual(t, "target not match", tc.name, tt.name)
		})
	}
}

func TestPriority_Remove(t *testing.T) {
	t.Parallel()
	t.Run("remove first", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 0, weight: 3, active: true}
		t1 := &Target{name: "t1", id: 1, weight: 2, active: true}
		t2 := &Target{name: "t2", id: 2, weight: 1, active: true}
		lb := zlb.NewPriority(t0, t1, t2)
		lb.Remove(0)
		tt, found := lb.Get(0)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t1", tt.name)
	})
	t.Run("remove middle", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 0, weight: 2, active: true}
		t1 := &Target{name: "t1", id: 1, weight: 3, active: true}
		t2 := &Target{name: "t2", id: 2, weight: 1, active: true}
		lb := zlb.NewPriority(t0, t1, t2)
		lb.Remove(1)
		tt, found := lb.Get(0)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t0", tt.name)
	})
	t.Run("remove last", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 0, weight: 1, active: true}
		t1 := &Target{name: "t1", id: 1, weight: 2, active: true}
		t2 := &Target{name: "t2", id: 2, weight: 3, active: true}
		lb := zlb.NewPriority(t0, t1, t2)
		lb.Remove(2)
		tt, found := lb.Get(0)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t1", tt.name)
	})
}
