package zos

import (
	"testing"

	"github.com/aileron-projects/go/ztesting"
)

func TestEnvSubst(t *testing.T) {
	t.Setenv("TestEnvSubst", "test")
	t.Run("resolve all", func(t *testing.T) {
		txt := `${TestEnvSubst}`
		b, err := EnvSubst([]byte(txt))
		ztesting.AssertEqual(t, "value not match", "test", string(b))
		ztesting.AssertEqualErr(t, "error not nil", nil, err)
	})
	t.Run("env error", func(t *testing.T) {
		txt := `${!TestEnvSubst}`
		b, err := EnvSubst([]byte(txt))
		ztesting.AssertEqual(t, "value not match", "", string(b))
		ztesting.AssertEqualErr(t, "error not match", &EnvError{Type: typeSyntax}, err)
	})
}

func TestEnvSubst2(t *testing.T) {
	t.Setenv("TestEnvSubst2", "test")
	t.Setenv("TestEnvSubst2_PTR", "TestEnvSubst2")
	t.Run("resolve all", func(t *testing.T) {
		txt := `${${TestEnvSubst2_PTR}}`
		b, err := EnvSubst2([]byte(txt))
		ztesting.AssertEqual(t, "value not match", "test", string(b))
		ztesting.AssertEqualErr(t, "error not nil", nil, err)
	})
	t.Run("inner env error", func(t *testing.T) {
		txt := `${${!TestEnvSubst2_PTR}}`
		b, err := EnvSubst2([]byte(txt))
		ztesting.AssertEqual(t, "value not match", "", string(b))
		ztesting.AssertEqualErr(t, "error not match", &EnvError{Type: typeSyntax}, err)
	})
	t.Run("outer env error", func(t *testing.T) {
		txt := `${!${TestEnvSubst2_PTR}}`
		b, err := EnvSubst2([]byte(txt))
		ztesting.AssertEqual(t, "value not match", "", string(b))
		ztesting.AssertEqualErr(t, "error not match", &EnvError{Type: typeSyntax}, err)
	})
}

func TestResolveEnv(t *testing.T) {
	t.Setenv("TestResolveEnv", "test")
	t.Setenv("TestResolveEnv_Cap", "TEST")
	testCases := map[string]struct {
		in     string
		result string
		err    error
	}{
		"case01": {"", "", &EnvError{Type: typeSyntax}},
		"case02": {"BAD", "", &EnvError{Type: typeSyntax}},
		"case03": {"${!prefix}", "", &EnvError{Type: typeSyntax}},
		"case04": {"${/*pattern/string}", "", &EnvError{Type: typeSyntax}},
		"case05": {"${parameter*other}", "", &EnvError{Type: typeSyntax}},
		"case06": {"${!TestResolve*}", "TestResolveEnv TestResolveEnv_Cap", nil},
		"case07": {"${!TestResolve@}", "TestResolveEnv TestResolveEnv_Cap", nil},
		"case08": {"${#TestResolveEnv}", "4", nil},
		"case09": {"${TestResolveEnv_UnDef:-word}", "word", nil},
		"case10": {"${TestResolveEnv_UnDef-word}", "word", nil},
		"case11": {"${TestResolveEnv:=word}", "test", nil},
		"case12": {"${TestResolveEnv=word}", "test", nil},
		"case13": {"${TestResolveEnv:?word}", "test", nil},
		"case14": {"${TestResolveEnv?word}", "test", nil},
		"case15": {"${TestResolveEnv:+word}", "word", nil},
		"case16": {"${TestResolveEnv+word}", "word", nil},
		"case17": {"${TestResolveEnv:2}", "st", nil},
		"case18": {"${TestResolveEnv:1:2}", "es", nil},
		"case19": {"${TestResolveEnv#t}", "est", nil},
		"case20": {"${TestResolveEnv##t}", "est", nil},
		"case21": {"${TestResolveEnv%t}", "tes", nil},
		"case22": {"${TestResolveEnv%%t}", "tes", nil},
		"case23": {"${TestResolveEnv/[t]/T}", "Test", nil},
		"case24": {"${TestResolveEnv//[^t]/X}", "tXXt", nil},
		"case25": {"${TestResolveEnv_Cap/#[T]/t}", "tEST", nil},
		"case26": {"${TestResolveEnv_Cap/%[T]/t}", "TESt", nil},
		"case27": {"${TestResolveEnv^[t]}", "Test", nil},
		"case28": {"${TestResolveEnv^^[^t]}", "tESt", nil},
		"case29": {"${TestResolveEnv_Cap,[T]}", "tEST", nil},
		"case30": {"${TestResolveEnv_Cap,,[^T]}", "TesT", nil},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := ResolveEnv([]byte(tc.in))
			ztesting.AssertEqual(t, "value not match", tc.result, string(v))
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
		})
	}
}

func TestResolveGroup1(t *testing.T) {
	t.Setenv("TestResolveGroup1", "test")
	testCases := map[string]struct {
		o      string
		result string
		err    error
	}{
		"case01": {"!TestResolve*", "TestResolveGroup1", nil},
		"case02": {"!TestResolve@", "TestResolveGroup1", nil},
		"case03": {"#TestResolveGroup1", "4", nil},
		"case04": {"", "", errSyntax("undefined", nil)},
		"case05": {"BAD", "", errSyntax("undefined", nil)},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := resolveGroup1(tc.o)
			ztesting.AssertEqual(t, "value not match", tc.result, v)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
		})
	}
}

func TestResolveGroup2(t *testing.T) {
	t.Setenv("TestResolveGroup2", "12345678")
	testCases := map[string]struct {
		p, o   string
		result string
		err    error
	}{
		"case01": {"TestResolveGroup2_UnDef", ":-word", "word", nil},
		"case02": {"TestResolveGroup2_UnDef", "-word", "word", nil},
		"case03": {"TestResolveGroup2", ":=word", "12345678", nil},
		"case04": {"TestResolveGroup2", "=word", "12345678", nil},
		"case05": {"TestResolveGroup2", ":?word", "12345678", nil},
		"case06": {"TestResolveGroup2", "?word", "12345678", nil},
		"case07": {"TestResolveGroup2", ":+word", "word", nil},
		"case08": {"TestResolveGroup2", "+word", "word", nil},
		"case09": {"TestResolveGroup2", ":5", "678", nil},
		"case10": {"TestResolveGroup2", ":3:2", "45", nil},
		"case11": {"TestResolveGroup2", "", "", errSyntax("undefined", nil)},
		"case12": {"TestResolveGroup2", "!BAD", "", errSyntax("undefined", nil)},
		"case13": {"TestResolveGroup2", ":", "", errSyntax("undefined", nil)},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := resolveGroup2(tc.p, tc.o)
			ztesting.AssertEqual(t, "value not match", tc.result, v)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
		})
	}
}

func TestResolveGroup3(t *testing.T) {
	t.Setenv("TestResolveGroup3", "test")
	testCases := map[string]struct {
		p, o   string
		result string
		err    error
	}{
		"case01": {"TestResolveGroup3", "#t", "est", nil},
		"case02": {"TestResolveGroup3", "##t", "est", nil},
		"case03": {"TestResolveGroup3", "%t", "tes", nil},
		"case04": {"TestResolveGroup3", "%%t", "tes", nil},
		"case05": {"TestResolveGroup3", "", "", errSyntax("undefined", nil)},
		"case06": {"TestResolveGroup3", "!BAD", "", errSyntax("undefined", nil)},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := resolveGroup3(tc.p, tc.o)
			ztesting.AssertEqual(t, "value not match", tc.result, v)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
		})
	}
}

func TestResolveGroup4(t *testing.T) {
	t.Setenv("TestResolveGroup4_01", "test")
	t.Setenv("TestResolveGroup4_02", "TEST")
	testCases := map[string]struct {
		p, o   string
		result string
		err    error
	}{
		"case01": {"TestResolveGroup4_01", "/[t]/T", "Test", nil},
		"case02": {"TestResolveGroup4_01", "//[^t]/X", "tXXt", nil},
		"case03": {"TestResolveGroup4_02", "/#[T]/t", "tEST", nil},
		"case04": {"TestResolveGroup4_02", "/%[T]/t", "TESt", nil},
		"case05": {"TestResolveGroup4", "", "", errSyntax("undefined", nil)},
		"case06": {"TestResolveGroup4", "/!BAD", "", errSyntax("undefined", nil)},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := resolveGroup4(tc.p, tc.o)
			ztesting.AssertEqual(t, "value not match", tc.result, v)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
		})
	}
}

func TestResolveGroup5(t *testing.T) {
	t.Setenv("TestResolveGroup5_01", "test")
	t.Setenv("TestResolveGroup5_02", "TEST")
	testCases := map[string]struct {
		p, o   string
		result string
		err    error
	}{
		"case01": {"TestResolveGroup5_01", "^[t]", "Test", nil},
		"case02": {"TestResolveGroup5_01", "^^[^t]", "tESt", nil},
		"case03": {"TestResolveGroup5_02", ",[T]", "tEST", nil},
		"case04": {"TestResolveGroup5_02", ",,[^T]", "TesT", nil},
		"case05": {"TestResolveGroup5", "", "", errSyntax("undefined", nil)},
		"case06": {"TestResolveGroup5", "!BAD", "", errSyntax("undefined", nil)},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			v, err := resolveGroup5(tc.p, tc.o)
			ztesting.AssertEqual(t, "value not match", tc.result, v)
			ztesting.AssertEqualErr(t, "error not match", tc.err, err)
		})
	}
}
