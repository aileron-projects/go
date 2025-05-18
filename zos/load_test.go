package zos_test

import (
	"testing"

	"github.com/aileron-projects/go/zos"
	"github.com/aileron-projects/go/ztesting"
)

func TestLoaEnv(t *testing.T) {
	t.Parallel()
	t.Run("comment", func(t *testing.T) {
		txt := `
		# comment line
		FOO=foo # inline comment
		`
		m, err := zos.LoadEnv([]byte(txt))
		want := map[string]string{"FOO": "foo"}
		for k, v := range want {
			ztesting.AssertEqual(t, "value not match", v, m[k])
		}
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	})
	t.Run("export", func(t *testing.T) {
		txt := `
		export FOO=foo
		export B_A_R=bar
		`
		m, err := zos.LoadEnv([]byte(txt))
		want := map[string]string{"FOO": "foo", "B_A_R": "bar"}
		for k, v := range want {
			ztesting.AssertEqual(t, "value not match", v, m[k])
		}
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	})
	t.Run("char escape", func(t *testing.T) {
		txt := `FOO=\f\o\o`
		m, err := zos.LoadEnv([]byte(txt))
		want := map[string]string{"FOO": "foo"}
		for k, v := range want {
			ztesting.AssertEqual(t, "value not match", v, m[k])
		}
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	})
	t.Run("quotations", func(t *testing.T) {
		txt := `
		NONE=none
		SINGLE='single'
		DOUBLE="double"
		`
		m, err := zos.LoadEnv([]byte(txt))
		want := map[string]string{"NONE": "none", "SINGLE": "single", "DOUBLE": "double"}
		for k, v := range want {
			ztesting.AssertEqual(t, "value not match", v, m[k])
		}
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	})
	t.Run("quotations in quotations", func(t *testing.T) {
		txt := `
		SINGLE='single and "double"'
		DOUBLE="'single' and double"
		`
		m, err := zos.LoadEnv([]byte(txt))
		want := map[string]string{"SINGLE": "single and \"double\"", "DOUBLE": "'single' and double"}
		for k, v := range want {
			ztesting.AssertEqual(t, "value not match", v, m[k])
		}
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	})
	t.Run("quotations escape", func(t *testing.T) {
		txt := `
		SINGLE='single \'escape\''
		DOUBLE="double \"escape\""
		`
		m, err := zos.LoadEnv([]byte(txt))
		want := map[string]string{"SINGLE": "single 'escape'", "DOUBLE": "double \"escape\""}
		for k, v := range want {
			ztesting.AssertEqual(t, "value not match", v, m[k])
		}
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	})
	t.Run("quotations sequence", func(t *testing.T) {
		txt := `
		SEQ1='Single'"Double"
		SEQ2='Single' and "Double"
		`
		m, err := zos.LoadEnv([]byte(txt))
		want := map[string]string{"SEQ1": "SingleDouble", "SEQ2": "Single and Double"}
		for k, v := range want {
			ztesting.AssertEqual(t, "value not match", v, m[k])
		}
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	})
	t.Run("multiline", func(t *testing.T) {
		txt := `
		MULTI1='
		line1
		line2
		'
		MULTI2="
		line1
		line2
		"
		`
		m, err := zos.LoadEnv([]byte(txt))
		want := map[string]string{"MULTI1": "line1line2", "MULTI2": "line1line2"}
		for k, v := range want {
			ztesting.AssertEqual(t, "value not match", v, m[k])
		}
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	})
	t.Run("multiline with line break", func(t *testing.T) {
		txt := `
		MULTI1='
		line1\n
		line2
		'
		MULTI2="
		line1\n
		line2
		"
		`
		m, err := zos.LoadEnv([]byte(txt))
		want := map[string]string{"MULTI1": "line1\nline2", "MULTI2": "line1\nline2"}
		for k, v := range want {
			ztesting.AssertEqual(t, "value not match", v, m[k])
		}
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	})
	t.Run("end with escape", func(t *testing.T) {
		txt := `FOO=foo\`
		m, err := zos.LoadEnv([]byte(txt))
		want := map[string]string{"FOO": "foo\\"}
		for k, v := range want {
			ztesting.AssertEqual(t, "value not match", v, m[k])
		}
		ztesting.AssertEqualErr(t, "non nil error returned", nil, err)
	})
	t.Run("env not found", func(t *testing.T) {
		txt := `=foo`
		m, err := zos.LoadEnv([]byte(txt))
		ztesting.AssertEqual(t, "unexpected value found", 0, len(m))
		ztesting.AssertEqualErr(t, "error not match", &zos.EnvError{Type: "zos: loading env failed."}, err)
	})
	t.Run("invalid char", func(t *testing.T) {
		txt := `***=foo`
		m, err := zos.LoadEnv([]byte(txt))
		ztesting.AssertEqual(t, "unexpected value found", 0, len(m))
		ztesting.AssertEqualErr(t, "error not match", &zos.EnvError{Type: "zos: loading env failed."}, err)
	})
	t.Run("invalid line format", func(t *testing.T) {
		txt := `foo`
		m, err := zos.LoadEnv([]byte(txt))
		ztesting.AssertEqual(t, "unexpected value found", 0, len(m))
		ztesting.AssertEqualErr(t, "error not match", &zos.EnvError{Type: "zos: loading env failed."}, err)
	})
	t.Run("quotation not closed", func(t *testing.T) {
		txt := `
		MULTI='
		line1
		line2
		`
		m, err := zos.LoadEnv([]byte(txt))
		ztesting.AssertEqual(t, "unexpected value found", 0, len(m))
		ztesting.AssertEqualErr(t, "error not match", &zos.EnvError{Type: "zos: loading env failed."}, err)
	})
	t.Run("env subst error", func(t *testing.T) {
		txt := `FOO=${!FOO}`
		m, err := zos.LoadEnv([]byte(txt))
		ztesting.AssertEqual(t, "unexpected value found", 0, len(m))
		ztesting.AssertEqualErr(t, "error not match", &zos.EnvError{Type: "zos: loading env failed."}, err)
	})
}
