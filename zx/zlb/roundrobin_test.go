package zlb_test

import (
	"testing"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/zx/zlb"
)

var rrTargets = map[string]struct {
	found   bool
	names   []string
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
		names: []string{"t0", "t0", "t0"},
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
		names: []string{"t1", "t0", "t1", "t1", "t0", "t1", "t1", "t0"},
		targets: []*Target{
			{name: "t0", id: 0, weight: 1, active: true},
			{name: "t1", id: 1, weight: 2, active: true},
		},
	},
	"2 targets, equal weight": {
		found: true,
		names: []string{"t0", "t1", "t0", "t1", "t0", "t1"},
		targets: []*Target{
			{name: "t0", id: 0, weight: 2, active: true},
			{name: "t1", id: 1, weight: 2, active: true},
		},
	},
	"3 targets, weights 321": {
		found: true,
		names: []string{"t0", "t1", "t0", "t2", "t1", "t0"},
		targets: []*Target{
			{name: "t0", id: 0, weight: 3, active: true},
			{name: "t1", id: 1, weight: 2, active: true},
			{name: "t2", id: 2, weight: 1, active: true},
		},
	},
	"3 targets, weights 231": {
		found: true,
		names: []string{"t1", "t0", "t1", "t2", "t0", "t1"},
		targets: []*Target{
			{name: "t0", id: 0, weight: 2, active: true},
			{name: "t1", id: 1, weight: 3, active: true},
			{name: "t2", id: 2, weight: 1, active: true},
		},
	},
	"3 targets, weights 123": {
		found: true,
		names: []string{"t2", "t1", "t0", "t2", "t1", "t2"},
		targets: []*Target{
			{name: "t0", id: 0, weight: 1, active: true},
			{name: "t1", id: 1, weight: 2, active: true},
			{name: "t2", id: 2, weight: 3, active: true},
		},
	},
	"3 targets, contains inactive": {
		found: true,
		names: []string{"t1", "t0", "t1", "t1", "t0", "t1"},
		targets: []*Target{
			{name: "t0", id: 0, weight: 1, active: true},
			{name: "t1", id: 1, weight: 2, active: true},
			{name: "t2", id: 2, weight: 3, active: false},
		},
	},
}

func TestRoundRobin_Get(t *testing.T) {
	t.Parallel()
	for name, tc := range rrTargets {
		t.Run(name, func(t *testing.T) {
			lb := zlb.NewRoundRobin(tc.targets...)
			if !tc.found {
				_, found := lb.Get(0)
				ztesting.AssertEqual(t, "found not match", tc.found, found)
				return
			}
			history := []string{}
			for range tc.names {
				tt, found := lb.Get(0)
				history = append(history, tt.name)
				ztesting.AssertEqual(t, "found not match", tc.found, found)
				ztesting.AssertEqual(t, "found not match", true, found)
			}
			ztesting.AssertEqualSlice(t, "history not match", tc.names, history)
		})
	}
}

func TestRoundRobin_Remove(t *testing.T) {
	t.Parallel()
	t.Run("remove first", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 0, weight: 3, active: true}
		t1 := &Target{name: "t1", id: 1, weight: 2, active: true}
		t2 := &Target{name: "t2", id: 2, weight: 1, active: true}
		lb := zlb.NewRoundRobin(t0, t1, t2)
		lb.Remove(0)
		ztesting.AssertEqualSlice(t, "targets not match", []*Target{t1, t2}, lb.Targets())
		history := []string{}
		want := []string{"t1", "t2", "t1", "t1", "t2"}
		for range want {
			tt, found := lb.Get(0)
			history = append(history, tt.name)
			ztesting.AssertEqual(t, "found not match", true, found)
		}
		ztesting.AssertEqualSlice(t, "history not match", want, history)
	})
	t.Run("remove middle", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 0, weight: 2, active: true}
		t1 := &Target{name: "t1", id: 1, weight: 3, active: true}
		t2 := &Target{name: "t2", id: 2, weight: 1, active: true}
		lb := zlb.NewRoundRobin(t0, t1, t2)
		lb.Remove(1)
		ztesting.AssertEqualSlice(t, "targets not match", []*Target{t0, t2}, lb.Targets())
		history := []string{}
		want := []string{"t0", "t2", "t0", "t0", "t2"}
		for range want {
			tt, found := lb.Get(0)
			history = append(history, tt.name)
			ztesting.AssertEqual(t, "found not match", true, found)
		}
		ztesting.AssertEqualSlice(t, "history not match", want, history)
	})
	t.Run("remove last", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 0, weight: 1, active: true}
		t1 := &Target{name: "t1", id: 1, weight: 2, active: true}
		t2 := &Target{name: "t2", id: 2, weight: 3, active: true}
		lb := zlb.NewRoundRobin(t0, t1, t2)
		lb.Remove(2)
		ztesting.AssertEqualSlice(t, "targets not match", []*Target{t0, t1}, lb.Targets())
		history := []string{}
		want := []string{"t1", "t0", "t1", "t1", "t0"}
		for range want {
			tt, found := lb.Get(0)
			history = append(history, tt.name)
			ztesting.AssertEqual(t, "found not match", true, found)
		}
		ztesting.AssertEqualSlice(t, "history not match", want, history)
	})
	t.Run("remove multiple", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 0, weight: 1, active: true}
		t1 := &Target{name: "t1", id: 1, weight: 2, active: true}
		t2 := &Target{name: "t2", id: 2, weight: 3, active: true}
		lb := zlb.NewRoundRobin(t2, t0, t2, t2, t2, t1, t2, t2)
		lb.Remove(2)
		ztesting.AssertEqualSlice(t, "targets not match", []*Target{t0, t1}, lb.Targets())
		history := []string{}
		want := []string{"t1", "t0", "t1", "t1", "t0"}
		for range want {
			tt, found := lb.Get(0)
			history = append(history, tt.name)
			ztesting.AssertEqual(t, "found not match", true, found)
		}
		ztesting.AssertEqualSlice(t, "history not match", want, history)
	})
}
