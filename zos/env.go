package zos

import (
	"os"
	"regexp"
	"strings"
)

// GetenvSlice get slice data from environmental variable.
// delim is the delimiter that separates each values.
func GetenvSlice(name, delim string) []string {
	return strings.Split(os.Getenv(name), delim)
}

// GetenvMap get map data from environmental variable.
// delim is the delimiter that separates key value pairs.
// sep is the separator string that separates key and value.
// For example ENV="foo=f1,foo=f2,bar=b1,baz"
// has ',' as delimiter and '=' as separator.
// It results in {"foo":["f1","f2"], "bar":["b1"], "baz":[""]}
// Key-value pairs without separator are considered as key only.
// In that cases, key is saved in the returned map with empty string value.
func GetenvMap(name, delim, sep string) map[string][]string {
	v := os.Getenv(name)
	vv := strings.Split(v, delim)
	m := make(map[string][]string, len(vv))
	for _, kv := range vv {
		before, after, _ := strings.Cut(kv, sep)
		m[before] = append(m[before], after)
	}
	return m
}

var (
	// envExp is the regular expression that matches to environmental variables.
	// See the [ResolveEnv] for considered patterns.
	envExp = `\$\{(` +
		`[0-9a-zA-Z_]+[@]?[UuLl]?` + `|` +
		`[#!][0-9a-zA-Z_]+[*@]?` + `|` +
		`[0-9a-zA-Z_]+[:\-=?+#%/,^][^\$]*` +
		`)\}`
	envRe = regexp.MustCompilePOSIX(envExp)
)

// EnvSubst substitute environmental variable in the given bytes.
// See the [ResolveEnv] for available variable syntax.
// EnvSubst does not support nested variables.
// Use [EnvSubst2] to allow 2 levels nested variable.
// Note that escaping variable like '\${FOO}' is not supported.
func EnvSubst(b []byte) (subst []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			subst = nil
			err, _ = r.(error)
		}
	}()
	subst = envRe.ReplaceAllFunc(b, func(b []byte) []byte {
		b, err = ResolveEnv(b)
		if err != nil {
			panic(err)
		}
		return b
	})
	return subst, nil
}

// EnvSubst2 substitute environmental variable in the given bytes.
// See the [ResolveEnv] for available variable syntax.
// EnvSubst does support nested variables up to 2 levels.
// ${FOO_${BAR}} is allowed but ${FOO_${BAR_${BAZ}}}} is not allowed.
// Note that escaping variable like '\${FOO}' is not supported.
func EnvSubst2(b []byte) (subst []byte, err error) {
	defer func() {
		r := recover()
		if r != nil {
			subst = nil
			err, _ = r.(error)
		}
	}()
	subst = envRe.ReplaceAllFunc(b, func(b []byte) []byte {
		b, err = ResolveEnv(b)
		if err != nil {
			panic(err)
		}
		return b
	})
	subst = envRe.ReplaceAllFunc(subst, func(b []byte) []byte {
		b, err = ResolveEnv(b)
		if err != nil {
			panic(err)
		}
		return b
	})
	return subst, nil
}

// ResolveEnv substitutes a single environmental variable expression.
// Supported expressions are listed below.
// Expressions are basically derived from shell parameter substitution.
// Note that the substitution behavior is NOT exactly the same as bash.
//   - https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html
//   - https://tldp.org/LDP/abs/html/parameter-substitution.html
//   - https://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_06_02
//
// Rules:
//
//	Expressions:
//	  01: ${parameter}                  --- See the substitution rule table below.
//	  02: ${parameter:-word}            --- See the substitution rule table below.
//	  03: ${parameter-word}             --- See the substitution rule table below.
//	  04: ${parameter:=word}            --- See the substitution rule table below.
//	  05: ${parameter=word}             --- See the substitution rule table below.
//	  06: ${parameter:?word}            --- See the substitution rule table below.
//	  07: ${parameter?word}             --- See the substitution rule table below.
//	  08: ${parameter:+word}            --- See the substitution rule table below.
//	  09: ${parameter+word}             --- See the substitution rule table below.
//	  10: ${parameter:offset}           --- Trim characters before offset.
//	  11: ${parameter:offset:length}    --- Trim characters before offset and after offset+length.
//	  12: ${!prefix*}                   --- Join the parameter name which has the prefix with a white space (Same with ${!prefix*}).
//	  13: ${!prefix@}                   --- Currently fallback to #12.
//	  14: ${#parameter}                 --- Length of value.
//	  15: ${parameter#word}             --- Currently fallback to #16.
//	  16: ${parameter##word}            --- Remove prefix of the value which matched to the word. Longest match if pattern specified.
//	  17: ${parameter%word}             --- Currently fallback to #18.
//	  18: ${parameter%%word}            --- Remove suffix of the value which matched to the word. Longest match if pattern specified.
//	  19: ${parameter/pattern/string}   --- Replace the first value which matched to the pattern to string.
//	  20: ${parameter//pattern/string}  --- Replace all values which matched to the pattern to string.
//	  21: ${parameter/#pattern/string}  --- Replace the prefix to string if matched to the pattern.
//	  22: ${parameter/%pattern/string}  --- Replace the suffix to string if matched to the pattern.
//	  23: ${parameter^pattern}          --- Convert initial character to upper case if matched to the pattern.
//	  24: ${parameter^^pattern}         --- Convert all characters which matched to the pattern to upper case.
//	  25: ${parameter,pattern}          --- Convert initial character to lower case if matched to the pattern.
//	  26: ${parameter,,pattern}         --- Convert all characters which matched to the pattern to lower case.
//	  27: ${parameter@operator}         --- Process value with the operator.
//
//	Substitution rules:
//	  |  #  |     expression     |    parameter Set     |  parameter Set  | parameter Unset |
//	  |     |                    |    and Not Null      |    But Null     |                 |
//	  | --- | ------------------ | -------------------- | --------------- | --------------- |
//	  | 01  | ${parameter}       | substitute parameter | substitute null | substitute null |
//	  | 02  | ${parameter:-word} | substitute parameter | substitute word | substitute word |
//	  | 03  | ${parameter-word}  | substitute parameter | substitute null | substitute word |
//	  | 04  | ${parameter:=word} | substitute parameter | substitute word | assign word     |
//	  | 05  | ${parameter=word}  | substitute parameter | substitute null | assign word     |
//	  | 06  | ${parameter:?word} | substitute parameter | error           | error           |
//	  | 07  | ${parameter?word}  | substitute parameter | substitute null | error           |
//	  | 08  | ${parameter:+word} | substitute word      | substitute null | substitute null |
//	  | 09  | ${parameter+word}  | substitute word      | substitute word | substitute null |
//
//	parameter:
//	  [0-9a-zA-Z_]+
//
//	word:
//	  [^\$]*
//
//	pattern:
//	  c       : matches to the character ('$' is not allowed).
//	  [a-z]   : matches specified character range.
//	  .*      : matches any length of characters.
//	  .?      : matches zero or single characters.
//
//	operator:
//	  U       : convert all characters to upper case using [strings.ToUpper]
//	  u       : convert the first character to upper case using [strings.ToUpper]
//	  L       : convert all characters to lower case using [strings.ToLower]
//	  l       : convert the first character to lower case using [strings.ToLower]
func ResolveEnv(in []byte) ([]byte, error) {
	if len(in) < 3 || string(in[:2]) != "${" || in[len(in)-1] != '}' {
		return nil, &EnvError{Type: typeSyntax, Info: "resolving `" + string(in) + "`"}
	}

	parameter, others := splitVar(in[2 : len(in)-1])
	var value string
	var err error
	if len(others) == 0 {
		// Pattern:
		//  ${parameter}
		value, _ = env01(string(parameter))
		return []byte(value), nil
	}

	if len(parameter) == 0 {
		// Pattern:
		//  ${!prefix*} ${!prefix@}
		//  ${#parameter}
		value, err = resolveGroup1(string(others))
		if err != nil {
			return nil, &EnvError{Err: err, Type: typeSyntax, Info: "resolving `" + string(in) + "`"}
		}
		return []byte(value), nil
	}

	switch others[0] {
	case '-', '=', '?', '+':
		//  ${parameter-word} ${parameter=word}
		//  ${parameter?word} ${parameter+word}
		value, err = resolveGroup2(string(parameter), string(others))
	case ':':
		// Pattern:
		//  ${parameter:-word} ${parameter:=word}
		//  ${parameter:?word} ${parameter:+word}
		//  ${parameter:offset} ${parameter:offset:length}
		value, err = resolveGroup2(string(parameter), string(others))
	case '#', '%':
		// Pattern:
		//  ${parameter#word} ${parameter##word}
		//  ${parameter%word} ${parameter%%word}
		value, err = resolveGroup3(string(parameter), string(others))
	case '/':
		// Pattern:
		//  ${parameter/pattern/string} ${parameter//pattern/string}
		//  ${parameter/#pattern/string} ${parameter/%pattern/string}
		value, err = resolveGroup4(string(parameter), string(others))
	case '^', ',':
		// Pattern:
		//  ${parameter^pattern} ${parameter^^pattern}
		//  ${parameter,pattern} ${parameter,,pattern}
		value, err = resolveGroup5(string(parameter), string(others))
	case '@':
		// Pattern:
		//  ${parameter@operator}
		value, err = env27(string(parameter), string(others[1:]))
	default:
		err = errSyntax("undefined", nil)
	}
	if err != nil {
		return nil, &EnvError{Err: err, Type: typeSyntax, Info: "resolving `" + string(in) + "`"}
	}
	return []byte(value), nil
}

func validEnvName(s string) bool {
	for i := 0; i < len(s); i++ {
		if !validEnvChar(s[i]) {
			return false
		}
	}
	return true
}

func validEnvChar(c byte) bool {
	switch {
	case '0' <= c && c <= '9':
		return true
	case 'a' <= c && c <= 'z':
		return true
	case 'A' <= c && c <= 'Z':
		return true
	case c == '_':
		return true
	default:
		return false
	}
}

func splitVar(b []byte) (parameter, others []byte) {
	for i, c := range b {
		if !validEnvChar(c) {
			return b[:i], b[i:]
		}
	}
	return b, nil
}

// resolveGroup1 resolves the following pattern group.
//   - 12: ${!prefix*}   --> o=!prefix*
//   - 13: ${!prefix@}   --> o=!prefix@
//   - 14: ${#parameter} --> o=#parameter
func resolveGroup1(o string) (string, error) {
	if len(o) < 2 {
		return "", errSyntax("undefined", nil)
	}
	switch o[0] {
	case '#':
		return env14(o[1:])
	case '!':
		switch o[len(o)-1] {
		case '*':
			return env12(o[1 : len(o)-1])
		case '@':
			return env13(o[1 : len(o)-1])
		}
	}
	return "", errSyntax("undefined", nil)
}

// resolveGroup2 resolves the following pattern group.
//   - 02: ${parameter:-word}         --> o=:-word
//   - 03: ${parameter-word}          --> o=-word
//   - 04: ${parameter:=word}         --> o=:=word
//   - 05: ${parameter=word}          --> o==word
//   - 06: ${parameter:?word}         --> o=:?word
//   - 07: ${parameter?word}          --> o=?word
//   - 08: ${parameter:+word}         --> o=:+word
//   - 09: ${parameter+word}          --> o=+word
//   - 10: ${parameter:offset}        --> o=:offset
//   - 11: ${parameter:offset:length} --> o=:offset:length
func resolveGroup2(p, o string) (string, error) {
	if len(o) < 1 {
		return "", errSyntax("undefined", nil)
	}
	switch o[0] {
	case '-':
		return env03(p, o[1:])
	case '=':
		return env05(p, o[1:])
	case '?':
		return env07(p, o[1:])
	case '+':
		return env09(p, o[1:])
	case ':':
		if len(o) < 2 {
			return "", errSyntax("undefined", nil)
		}
		switch o[1] {
		case '-':
			return env02(p, o[2:])
		case '=':
			return env04(p, o[2:])
		case '?':
			return env06(p, o[2:])
		case '+':
			return env08(p, o[2:])
		default:
			if i := strings.Index(o[1:], ":"); i > 0 {
				return env11(p, o[1:i+1], o[i+2:])
			} else {
				return env10(p, o[1:])
			}
		}
	}
	return "", errSyntax("undefined", nil)
}

// resolveGroup3 resolves the following pattern group.
//   - 15: ${parameter#word}   --> o=#word
//   - 16: ${parameter##word}  --> o=##word
//   - 17: ${parameter%word}   --> o=%word
//   - 18: ${parameter%%word}  --> o=%%word
func resolveGroup3(p, o string) (string, error) {
	if len(o) < 1 {
		return "", errSyntax("undefined", nil)
	}
	switch o[0] {
	case '#':
		if len(o) > 1 && o[1] == '#' {
			return env16(p, o[2:])
		} else {
			return env15(p, o[1:])
		}
	case '%':
		if len(o) > 1 && o[1] == '%' {
			return env18(p, o[2:])
		} else {
			return env17(p, o[1:])
		}
	}
	return "", errSyntax("undefined", nil)
}

// resolveGroup4 resolves the following pattern group.
//   - 19: ${parameter/pattern/string}  --> o=/pattern/string
//   - 20: ${parameter//pattern/string} --> o=//pattern/string
//   - 21: ${parameter/#pattern/string} --> o=/#pattern/string
//   - 22: ${parameter/%pattern/string} --> o=/%pattern/string
func resolveGroup4(p, o string) (string, error) {
	if len(o) < 2 {
		return "", errSyntax("undefined", nil)
	}
	switch o[1] {
	case '/':
		if i := strings.Index(o[2:], "/"); i > 0 {
			return env20(p, o[2:i+2], o[i+3:])
		}
	case '#':
		if i := strings.Index(o[2:], "/"); i > 0 {
			return env21(p, o[2:i+2], o[i+3:])
		}
	case '%':
		if i := strings.Index(o[2:], "/"); i > 0 {
			return env22(p, o[2:i+2], o[i+3:])
		}
	default:
		if i := strings.Index(o[1:], "/"); i > 0 {
			return env19(p, o[1:i+1], o[i+2:])
		}
	}
	return "", errSyntax("undefined", nil)
}

// resolveGroup5 resolves the following pattern group.
//   - 23: ${parameter^pattern}  --> o=^pattern
//   - 24: ${parameter^^pattern} --> o=^^pattern
//   - 25: ${parameter,pattern}  --> o=,pattern
//   - 26: ${parameter,,pattern} --> o=,,pattern
func resolveGroup5(p, o string) (string, error) {
	if len(o) < 1 {
		return "", errSyntax("undefined", nil)
	}
	switch o[0] {
	case '^':
		if len(o) > 1 && o[1] == '^' {
			return env24(p, o[2:])
		} else {
			return env23(p, o[1:])
		}
	case ',':
		if len(o) > 1 && o[1] == ',' {
			return env26(p, o[2:])
		} else {
			return env25(p, o[1:])
		}
	}
	return "", errSyntax("undefined", nil)
}
