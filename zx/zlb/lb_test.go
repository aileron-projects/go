package zlb

import (
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

type testTarget struct {
	name   string
	id     uint64
	weight uint16
	active bool
}

func (t *testTarget) ID() uint64 {
	return t.id
}

func (t *testTarget) Weight() uint16 {
	return t.weight
}

func (t *testTarget) Active() bool {
	return t.active
}

func TestBaseLB_Add(t *testing.T) {
	t.Parallel()

	t1 := &testTarget{name: "t1", id: 1}
	t2 := &testTarget{name: "t2", id: 2}
	t3 := &testTarget{name: "t3", id: 3}

	t.Run("add one", func(t *testing.T) {
		lb := &baseLB[*testTarget]{}
		lb.Add(t1)
		ts := lb.Targets()
		ztesting.AssertEqual(t, "length not match", 1, len(ts))
		ztesting.AssertEqual(t, "index 0 not match", "t1", ts[0].name)
	})
	t.Run("add multiple", func(t *testing.T) {
		lb := &baseLB[*testTarget]{}
		lb.Add(t1, t2, t3)
		ts := lb.Targets()
		ztesting.AssertEqual(t, "length not match", 3, len(ts))
		ztesting.AssertEqual(t, "index 0 not match", "t1", ts[0].name)
		ztesting.AssertEqual(t, "index 1 not match", "t2", ts[1].name)
		ztesting.AssertEqual(t, "index 2 not match", "t3", ts[2].name)
	})
}

func TestBaseLB_Remove(t *testing.T) {
	t.Parallel()

	t1 := &testTarget{name: "t1", id: 1}
	t2 := &testTarget{name: "t2", id: 2}
	t3 := &testTarget{name: "t3", id: 3}

	t.Run("mode0", func(t *testing.T) {
		t.Run("remove from empty", func(t *testing.T) {
			lb := &baseLB[*testTarget]{}
			lb.Remove(1)
			ztesting.AssertEqual(t, "length not match", 0, len(lb.targets))
		})
		t.Run("remove not match", func(t *testing.T) {
			lb := &baseLB[*testTarget]{}
			lb.Add(t1, t2, t3)
			lb.Remove(999)
			ztesting.AssertEqual(t, "length not match", 3, len(lb.targets))
		})
		t.Run("remove first", func(t *testing.T) {
			lb := &baseLB[*testTarget]{}
			lb.Add(t1, t2, t3)
			lb.Remove(1)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 2, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t2", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t3", ts[1].name)
		})
		t.Run("remove middle", func(t *testing.T) {
			lb := &baseLB[*testTarget]{}
			lb.Add(t1, t2, t3)
			lb.Remove(2)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 2, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t1", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t3", ts[1].name)
		})
		t.Run("remove last", func(t *testing.T) {
			lb := &baseLB[*testTarget]{}
			lb.Add(t1, t2, t3)
			lb.Remove(3)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 2, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t1", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t2", ts[1].name)
		})
		t.Run("remove multiple 1", func(t *testing.T) {
			lb := &baseLB[*testTarget]{}
			lb.Add(t1, t2, t3, t1)
			lb.Remove(1)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 2, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t2", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t3", ts[1].name)
		})
		t.Run("remove multiple 2", func(t *testing.T) {
			lb := &baseLB[*testTarget]{}
			lb.Add(t1, t2, t2, t3)
			lb.Remove(2)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 2, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t1", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t3", ts[1].name)
		})
		t.Run("remove complex", func(t *testing.T) {
			lb := &baseLB[*testTarget]{}
			lb.Add(t1, t2, t3, t2, t2, t1, t3, t2)
			lb.Remove(2)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 4, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t1", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t3", ts[1].name)
			ztesting.AssertEqual(t, "index 2 not match", "t1", ts[2].name)
			ztesting.AssertEqual(t, "index 3 not match", "t3", ts[3].name)
		})
	})

	t.Run("mode1", func(t *testing.T) {
		t.Run("remove from empty", func(t *testing.T) {
			lb := &baseLB[*testTarget]{removeMode: 1}
			lb.Remove(1)
			ztesting.AssertEqual(t, "length not match", 0, len(lb.targets))
		})
		t.Run("remove not match", func(t *testing.T) {
			lb := &baseLB[*testTarget]{removeMode: 1}
			lb.Add(t1, t2, t3)
			lb.Remove(999)
			ztesting.AssertEqual(t, "length not match", 3, len(lb.targets))
		})
		t.Run("remove first", func(t *testing.T) {
			lb := &baseLB[*testTarget]{removeMode: 1}
			lb.Add(t1, t2, t3)
			lb.Remove(1)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 2, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t3", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t2", ts[1].name)
		})
		t.Run("remove middle", func(t *testing.T) {
			lb := &baseLB[*testTarget]{removeMode: 1}
			lb.Add(t1, t2, t3)
			lb.Remove(2)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 2, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t1", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t3", ts[1].name)
		})
		t.Run("remove last", func(t *testing.T) {
			lb := &baseLB[*testTarget]{removeMode: 1}
			lb.Add(t1, t2, t3)
			lb.Remove(3)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 2, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t1", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t2", ts[1].name)
		})
		t.Run("remove multiple 1", func(t *testing.T) {
			lb := &baseLB[*testTarget]{removeMode: 1}
			lb.Add(t1, t2, t3, t1)
			lb.Remove(1)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 2, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t3", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t2", ts[1].name)
		})
		t.Run("remove multiple 2", func(t *testing.T) {
			lb := &baseLB[*testTarget]{removeMode: 1}
			lb.Add(t1, t2, t2, t3)
			lb.Remove(2)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 2, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t1", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t3", ts[1].name)
		})
		t.Run("remove complex", func(t *testing.T) {
			lb := &baseLB[*testTarget]{removeMode: 1}
			lb.Add(t1, t2, t3, t2, t2, t1, t3, t2)
			lb.Remove(2)
			ts := lb.Targets()
			ztesting.AssertEqual(t, "length not match", 4, len(ts))
			ztesting.AssertEqual(t, "index 0 not match", "t1", ts[0].name)
			ztesting.AssertEqual(t, "index 1 not match", "t3", ts[1].name)
			ztesting.AssertEqual(t, "index 2 not match", "t3", ts[2].name)
			ztesting.AssertEqual(t, "index 3 not match", "t1", ts[3].name)
		})
	})
}
