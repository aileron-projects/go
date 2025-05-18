package zlb_test

import (
	"math"
	"testing"

	"github.com/aileron-projects/go/ztesting"
	"github.com/aileron-projects/go/zx/zlb"
)

var ringTargets = map[string]struct {
	key     uint64
	found   bool
	name    string
	targets []*Target
}{
	"0 target": {
		key:   123,
		found: false,
	},
	"1 target, inactive": {
		key:   123,
		found: false,
		targets: []*Target{
			{name: "t0", id: 999000, weight: 1, active: false},
		},
	},
	"1 target, 0 weight": {
		key:   123,
		found: false,
		targets: []*Target{
			{name: "t0", id: 999000, weight: 0, active: true},
		},
	},
	"1 target, non-0 weight": {
		key:   123,
		found: true,
		name:  "t0",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 1, active: true},
		},
	},
	"2 targets, 0 weight": {
		key:   123,
		found: false,
		targets: []*Target{
			{name: "t0", id: 999000, weight: 0, active: true},
			{name: "t1", id: 999111, weight: 0, active: true},
		},
	},
	"2 targets, non-0 weight": {
		key:   123,
		found: true,
		name:  "t1",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 1, active: true},
			{name: "t1", id: 999111, weight: 2, active: true},
		},
	},
	"2 targets, equal weight": {
		key:   123,
		found: true,
		name:  "t0",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 2, active: true},
			{name: "t1", id: 999111, weight: 2, active: true},
		},
	},
	"3 targets, key12345": {
		key:   12345,
		found: true,
		name:  "t2",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 3, active: true},
			{name: "t1", id: 999111, weight: 2, active: true},
			{name: "t2", id: 999222, weight: 1, active: true},
		},
	},
	"3 targets, key67890": {
		key:   67890,
		found: true,
		name:  "t1",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 3, active: true},
			{name: "t1", id: 999111, weight: 2, active: true},
			{name: "t2", id: 999222, weight: 1, active: true},
		},
	},
	"3 targets, contains inactive": {
		key:   12345,
		found: true,
		name:  "t1",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 1, active: true},
			{name: "t1", id: 999111, weight: 2, active: true},
			{name: "t2", id: 999222, weight: 3, active: false},
		},
	},
	"max key": {
		key:   uint64(math.MaxUint64),
		found: true,
		name:  "t1",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 0, active: true},
			{name: "t1", id: 999111, weight: 1, active: true},
		},
	},
	"zero key": {
		key:   uint64(0),
		found: true,
		name:  "t1",
		targets: []*Target{
			{name: "t0", id: 999000, weight: 0, active: true},
			{name: "t1", id: 999111, weight: 1, active: true},
		},
	},
}

func TestRingHash_Get(t *testing.T) {
	t.Parallel()
	for name, tc := range ringTargets {
		t.Run(name, func(t *testing.T) {
			lb := zlb.NewRingHash(tc.targets...)
			tt, found := lb.Get(tc.key)
			ztesting.AssertEqual(t, "found not match", tc.found, found)
			if !tc.found {
				return
			}
			ztesting.AssertEqual(t, "target not match", tc.name, tt.name)
			ztesting.AssertEqual(t, "active status not match", true, tt.Active())
		})
	}
}

func TestRingHash_Remove(t *testing.T) {
	t.Parallel()
	t.Run("remove first", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 999000, weight: 3, active: true}
		t1 := &Target{name: "t1", id: 999111, weight: 2, active: true}
		t2 := &Target{name: "t2", id: 999222, weight: 1, active: true}
		lb := zlb.NewRingHash(t0, t1, t2)
		tt, found := lb.Get(123)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t0", tt.name)
		lb.Remove(999000)
		lb.Update()
		tt, found = lb.Get(123)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t2", tt.name)
	})
	t.Run("remove middle", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 999000, weight: 2, active: true}
		t1 := &Target{name: "t1", id: 999111, weight: 3, active: true}
		t2 := &Target{name: "t2", id: 999222, weight: 1, active: true}
		lb := zlb.NewRingHash(t0, t1, t2)
		tt, found := lb.Get(12345)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t1", tt.name)
		lb.Remove(999111)
		lb.Update()
		tt, found = lb.Get(12345)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t2", tt.name)
	})
	t.Run("remove last", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 999000, weight: 1, active: true}
		t1 := &Target{name: "t1", id: 999111, weight: 2, active: true}
		t2 := &Target{name: "t2", id: 999222, weight: 3, active: true}
		lb := zlb.NewRingHash(t0, t1, t2)
		tt, found := lb.Get(12345678)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t2", tt.name)
		lb.Remove(999222)
		lb.Update()
		tt, found = lb.Get(12345678)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t1", tt.name)
	})
}

func TestRingHash_Update(t *testing.T) {
	t.Parallel()
	t.Run("empty", func(t *testing.T) {
		lb := zlb.NewRingHash[*Target]()
		lb.Update()
		_, found := lb.Get(123)
		ztesting.AssertEqual(t, "found not match", false, found)
	})
	t.Run("same weight", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 999000, weight: 5, active: true}
		t1 := &Target{name: "t1", id: 999111, weight: 5, active: true}
		t2 := &Target{name: "t2", id: 999222, weight: 5, active: true}
		lb := zlb.NewRingHash(t0, t1, t2)
		lb.Update()
		tt, found := lb.Get(123)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t1", tt.name)
	})
	t.Run("same id", func(t *testing.T) {
		t0 := &Target{name: "t0", id: 123, weight: 2, active: true}
		t1 := &Target{name: "t1", id: 123, weight: 3, active: true}
		t2 := &Target{name: "t2", id: 123, weight: 1, active: true}
		lb := zlb.NewRingHash(t0, t1, t2)
		lb.Update()
		tt, found := lb.Get(123)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t0", tt.name)
		tt, found = lb.Get(456)
		ztesting.AssertEqual(t, "found not match", true, found)
		ztesting.AssertEqual(t, "target not match", "t0", tt.name)
	})
}
