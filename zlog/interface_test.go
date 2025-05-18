package zlog_test

import (
	"testing"

	"github.com/aileron-projects/go/zlog"
	"github.com/aileron-projects/go/ztesting"
)

func TestLevel_String(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		lv  zlog.Level
		str string
	}{
		"trace":        {zlog.LvTrace, "TRACE"},
		"debug":        {zlog.LvDebug, "DEBUG"},
		"info":         {zlog.LvInfo, "INFO"},
		"warn":         {zlog.LvWarn, "WARN"},
		"error":        {zlog.LvError, "ERROR"},
		"fatal":        {zlog.LvFatal, "FATAL"},
		"trace+1":      {zlog.LvTrace + 1, "TRACE"},
		"debug+1":      {zlog.LvDebug + 1, "DEBUG"},
		"info+1":       {zlog.LvInfo + 1, "INFO"},
		"warn+1":       {zlog.LvWarn + 1, "WARN"},
		"error+1":      {zlog.LvError + 1, "ERROR"},
		"fatal+1":      {zlog.LvFatal + 1, "FATAL"},
		"trace+2":      {zlog.LvTrace + 2, "TRACE"},
		"debug+2":      {zlog.LvDebug + 2, "DEBUG"},
		"info+2":       {zlog.LvInfo + 2, "INFO"},
		"warn+2":       {zlog.LvWarn + 2, "WARN"},
		"error+2":      {zlog.LvError + 2, "ERROR"},
		"fatal+2":      {zlog.LvFatal + 2, "FATAL"},
		"trace+3":      {zlog.LvTrace + 3, "TRACE"},
		"debug+3":      {zlog.LvDebug + 3, "DEBUG"},
		"info+3":       {zlog.LvInfo + 3, "INFO"},
		"warn+3":       {zlog.LvWarn + 3, "WARN"},
		"error+3":      {zlog.LvError + 3, "ERROR"},
		"fatal+3":      {zlog.LvFatal + 3, "FATAL"},
		"undefined 0":  {zlog.Level(0), "UNDEFINED"},
		"undefined 25": {zlog.Level(25), "UNDEFINED"},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ztesting.AssertEqual(t, "wrong string expression of log level.", tc.str, tc.lv.String())
		})
	}
}

func TestLevel(t *testing.T) {
	t.Parallel()

	t.Run("higher equal", func(t *testing.T) {
		ztesting.AssertEqual(t, "wrong comparison result.", false, zlog.Level(zlog.LvInfo-1).HigherEqual(zlog.LvInfo))
		ztesting.AssertEqual(t, "wrong comparison result.", true, zlog.LvInfo.HigherEqual(zlog.LvInfo))
		ztesting.AssertEqual(t, "wrong comparison result.", true, zlog.Level(zlog.LvInfo+1).HigherEqual(zlog.LvInfo))
	})

	t.Run("higher than", func(t *testing.T) {
		ztesting.AssertEqual(t, "wrong comparison result.", false, zlog.Level(zlog.LvInfo-1).HigherThan(zlog.LvInfo))
		ztesting.AssertEqual(t, "wrong comparison result.", false, zlog.LvInfo.HigherThan(zlog.LvInfo))
		ztesting.AssertEqual(t, "wrong comparison result.", true, zlog.Level(zlog.LvInfo+1).HigherThan(zlog.LvInfo))
	})

	t.Run("less equal", func(t *testing.T) {
		ztesting.AssertEqual(t, "wrong comparison result.", true, zlog.LvInfo.LessEqual(zlog.LvInfo+1))
		ztesting.AssertEqual(t, "wrong comparison result.", true, zlog.LvInfo.LessEqual(zlog.LvInfo))
		ztesting.AssertEqual(t, "wrong comparison result.", false, zlog.LvInfo.LessEqual(zlog.LvInfo-1))
	})

	t.Run("less than", func(t *testing.T) {
		ztesting.AssertEqual(t, "wrong comparison result.", true, zlog.LvInfo.LessThan(zlog.LvInfo+1))
		ztesting.AssertEqual(t, "wrong comparison result.", false, zlog.LvInfo.LessThan(zlog.LvInfo))
		ztesting.AssertEqual(t, "wrong comparison result.", false, zlog.LvInfo.LessThan(zlog.LvInfo-1))
	})
}
