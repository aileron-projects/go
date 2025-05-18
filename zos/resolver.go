package zos

import (
	"cmp"
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	typeInvalidName  = "zos: invalid env name. `[0-9a-zA-Z_]+` is allowed."
	typeSyntax       = "zos: invalid env syntax."
	typeSubstitution = "zos: env substitution."
	typeLoad         = "zos: loading env failed."
)

// EnvError is the environmental substitution error.
type EnvError struct {
	Err     error
	Type    string
	Pattern string
	Info    string
}

func (e *EnvError) Unwrap() error {
	return e.Err
}

func (e *EnvError) Error() string {
	s := e.Type
	if e.Pattern != "" {
		s += " (detected syntax is " + e.Pattern + ")"
	}
	s += " " + e.Info
	if e.Err != nil {
		return s + " [" + e.Err.Error() + "]"
	}
	return s
}

func (e *EnvError) Is(err error) bool {
	for err != nil {
		ee, ok := err.(*EnvError)
		if ok {
			return e.Type == ee.Type
		}
		err = errors.Unwrap(err)
	}
	return false
}

func errInvalidName(pattern, name string) *EnvError {
	return &EnvError{
		Type:    typeInvalidName,
		Pattern: pattern,
		Info:    "got " + name,
	}
}

func errSyntax(pattern string, err error) *EnvError {
	return &EnvError{
		Err:     err,
		Type:    typeSyntax,
		Pattern: pattern,
	}
}

func errSubstitute(pattern string, value string) *EnvError {
	return &EnvError{
		Type:    typeSubstitution,
		Pattern: pattern,
		Info:    "got " + value,
	}
}

// env01 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 01: ${parameter}
func env01(p string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter}", p)
	}
	return os.Getenv(p), nil
}

// env02 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 02: ${parameter:-word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: substitute word
//   - parameter Unset: substitute word
func env02(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter:-word}", p)
	}
	v, ok := os.LookupEnv(p)
	if !ok { // parameter unset
		return w, nil
	}
	if v == "" { // parameter set but null
		return w, nil
	}
	return v, nil // parameter set not null
}

// env03 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 03: ${parameter-word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: substitute null
//   - parameter Unset: substitute word
func env03(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter-word}", p)
	}
	v, ok := os.LookupEnv(p)
	if !ok { // parameter unset
		return w, nil
	}
	if v == "" { // parameter set but null
		return "", nil
	}
	return v, nil // parameter set not null
}

// env04 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 04: ${parameter:=word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: assign word
//   - parameter Unset: assign word
func env04(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter:=word}", p)
	}
	v, ok := os.LookupEnv(p)
	if !ok { // parameter unset
		return w, os.Setenv(p, w)
	}
	if v == "" { // parameter set but null
		return w, os.Setenv(p, w)
	}
	return v, nil // parameter set not null
}

// env05 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 05: ${parameter=word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: substitute null
//   - parameter Unset: assign word
func env05(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter=word}", p)
	}
	v, ok := os.LookupEnv(p)
	if !ok { // parameter unset
		return w, os.Setenv(p, w)
	}
	if v == "" { // parameter set but null
		return "", nil
	}
	return v, nil // parameter set not null
}

// env06 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 06: ${parameter:?word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: error
//   - parameter Unset: error
func env06(p, _ string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter:?word}", p)
	}
	v, ok := os.LookupEnv(p)
	if !ok { // parameter unset
		return "", errSubstitute("${parameter:?word}", p)
	}
	if v == "" { // parameter set but null
		return "", errSubstitute("${parameter:?word}", p)
	}
	return v, nil // parameter set not null
}

// env07 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 07: ${parameter?word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute parameter
//   - parameter Set But Null: substitute null
//   - parameter Unset: error
func env07(p, _ string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter?word}", p)
	}
	v, ok := os.LookupEnv(p)
	if !ok { // parameter unset
		return "", errSubstitute("${parameter?word}", p)
	}
	if v == "" { // parameter set but null
		return "", nil
	}
	return v, nil // parameter set not null
}

// env08 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 08: ${parameter:+word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute word
//   - parameter Set But Null: substitute null
//   - parameter Unset: substitute null
func env08(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter:+word}", p)
	}
	v, ok := os.LookupEnv(p)
	if !ok { // parameter unset
		return "", nil
	}
	if v == "" { // parameter set but null
		return "", nil
	}
	return w, nil // parameter set not null
}

// env09 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 09: ${parameter+word}
//
// This pattern works as:
//   - parameter Set and Not Null: substitute word
//   - parameter Set But Null: substitute word
//   - parameter Unset: substitute null
func env09(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter+word}", p)
	}
	v, ok := os.LookupEnv(p)
	if !ok { // parameter unset
		return "", nil
	}
	if v == "" { // parameter set but null
		return w, nil
	}
	return w, nil // parameter set not null
}

// env10 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 10: ${parameter:offset}
func env10(p, o string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter:offset}", p)
	}
	offset, err := strconv.Atoi(o)
	if err != nil {
		return "", errSyntax("${parameter:offset}", err)
	}
	v := os.Getenv(p)
	r := []rune(v)
	if offset < 0 {
		return v, nil
	}
	if offset > len(r) {
		return "", nil
	}
	return string(r[offset:]), nil
}

// env11 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 11: ${parameter:offset:length}
func env11(p, o, l string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter:offset:length}", p)
	}
	offset, err := strconv.Atoi(o)
	if err != nil {
		return "", errSyntax("${parameter:offset:length}", err)
	}
	length, err := strconv.Atoi(l)
	if err != nil {
		return "", errSyntax("${parameter:offset:length}", err)
	}
	v := os.Getenv(p)
	if offset < 0 {
		return v, nil
	}
	r := []rune(v)
	if offset >= len(r) {
		return "", nil
	}
	if length < 0 {
		return string(r[offset:]), nil
	}
	if offset+length > len(r) {
		return string(r[offset:]), nil
	}
	return string(r[offset : offset+length]), nil
}

// env12 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 12: ${!prefix*}
func env12(p string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${!prefix*}", p)
	}
	names := []string{}
	for _, v := range os.Environ() {
		if strings.HasPrefix(v, p) {
			names = append(names, strings.Split(v, "=")[0])
		}
	}
	return strings.Join(names, " "), nil
}

// env13 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 13: ${!prefix@}
func env13(p string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${!prefix@}", p)
	}
	return env12(p) // Fallback
}

// env14 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 14: ${#parameter}
func env14(p string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${#parameter}", p)
	}
	return strconv.Itoa(len([]rune(os.Getenv(p)))), nil
}

// env15 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 15: ${parameter#word}
func env15(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter#word}", p)
	}
	return env16(p, w) // Fallback.
}

// env16 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 16: ${parameter##word}
func env16(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter##word}", p)
	}
	re, err := regexp.CompilePOSIX("^" + w)
	if err != nil {
		return "", errSyntax("${parameter##word}", err)
	}
	v := os.Getenv(p)
	i := re.FindAllStringIndex(v, -1)
	if len(i) == 0 {
		return v, nil
	}
	return v[i[len(i)-1][1]:], nil
}

// env17 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 17: ${parameter%word}
func env17(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter%word}", p)
	}
	return env18(p, w) // Fallback.
}

// env18 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 18: ${parameter%%word}
func env18(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter%%word}", p)
	}
	re, err := regexp.CompilePOSIX(w + "$")
	if err != nil {
		return "", errSyntax("${parameter%%word}", err)
	}
	v := os.Getenv(p)
	i := re.FindAllStringIndex(v, -1)
	if len(i) == 0 {
		return v, nil
	}
	return v[:i[len(i)-1][0]], nil // Remove longest match.
}

// env19 returns the env value resolved from the following pattern.
// It replaces first matched pattern in env value to the string.
// See the [ResolveEnv].
//   - 19: ${parameter/pattern/string}
func env19(p, w, s string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter/pattern/string}", p)
	}
	re, err := regexp.CompilePOSIX(w)
	if err != nil {
		return "", errSyntax("${parameter/pattern/string}", err)
	}
	v := os.Getenv(p)
	replaced := false
	v = re.ReplaceAllStringFunc(v, func(ss string) string {
		if replaced {
			return ss
		}
		replaced = true
		return s
	})
	return v, nil
}

// env20 returns the env value resolved from the following pattern.
// It replaces all patterns matches to the input pattern to the string.
// See the [ResolveEnv].
//   - 20: ${parameter//pattern/string}
func env20(p, w, s string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter//pattern/string}", p)
	}
	re, err := regexp.CompilePOSIX(w)
	if err != nil {
		return "", errSyntax("${parameter//pattern/string}", err)
	}
	v := os.Getenv(p)
	return re.ReplaceAllString(v, s), nil
}

// env21 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 21: ${parameter/#pattern/string}
func env21(p, w, s string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter/#pattern/string}", p)
	}
	re, err := regexp.CompilePOSIX("^" + w)
	if err != nil {
		return "", errSyntax("${parameter/#pattern/string}", err)
	}
	v := os.Getenv(p)
	i := re.FindAllStringIndex(v, -1)
	if len(i) == 0 {
		return v, nil
	}
	return s + v[i[len(i)-1][1]:], nil
}

// env22 returns the env value resolved from the following pattern.
// See the [ResolveEnv].
//   - 22: ${parameter/%pattern/string}
func env22(p, w, s string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter/%pattern/string}", p)
	}
	re, err := regexp.CompilePOSIX(w + "$")
	if err != nil {
		return "", errSyntax("${parameter/%pattern/string}", err)
	}
	v := os.Getenv(p)
	i := re.FindAllStringIndex(v, -1)
	if len(i) == 0 {
		return v, nil
	}
	return v[:i[len(i)-1][0]] + s, nil
}

// env23 returns the env value resolved from the following pattern.
// Convert the first char of the value if matched to the pattern into upper.
// See the [ResolveEnv].
//   - 23: ${parameter^pattern}
func env23(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter^pattern}", p)
	}
	re, err := regexp.CompilePOSIX(cmp.Or(w, ".?"))
	if err != nil {
		return "", errSyntax("${parameter^pattern}", err)
	}
	v := os.Getenv(p)
	if len(v) == 0 {
		return "", nil
	}
	return re.ReplaceAllStringFunc(v[:1], strings.ToUpper) + v[1:], nil
}

// env24 returns the env value resolved from the following pattern.
// Convert all matched chars to upper.
// See the [ResolveEnv].
//   - 24: ${parameter^^pattern}
func env24(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter^^pattern}", p)
	}
	re, err := regexp.CompilePOSIX(cmp.Or(w, ".?"))
	if err != nil {
		return "", errSyntax("${parameter^^pattern}", err)
	}
	v := os.Getenv(p)
	return re.ReplaceAllStringFunc(v, strings.ToUpper), nil
}

// env25 returns the env value resolved from the following pattern.
// Convert the first char of the value if matched to the pattern into lower.
// See the [ResolveEnv].
//   - 25: ${parameter,pattern}
func env25(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter,pattern}", p)
	}
	re, err := regexp.CompilePOSIX(cmp.Or(w, ".?"))
	if err != nil {
		return "", errSyntax("${parameter,pattern}", err)
	}
	v := os.Getenv(p)
	if len(v) == 0 {
		return "", nil
	}
	return re.ReplaceAllStringFunc(v[:1], strings.ToLower) + v[1:], nil
}

// env26 returns the env value resolved from the following pattern.
// Convert all matched chars to lower.
// See the [ResolveEnv].
//   - 26: ${parameter,,pattern}
func env26(p, w string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter,,pattern}", p)
	}
	re, err := regexp.CompilePOSIX(cmp.Or(w, ".?"))
	if err != nil {
		return "", errSyntax("${parameter,,pattern}", err)
	}
	v := os.Getenv(p)
	return re.ReplaceAllStringFunc(v, strings.ToLower), nil
}

// env27 returns the env value resolved from the following pattern.
// Apply operation to the parameter.
// See the [ResolveEnv].
//   - 27: ${parameter@operator}
func env27(p, o string) (string, error) {
	if !validEnvName(p) {
		return "", errInvalidName("${parameter@operator}", p)
	}
	v := os.Getenv(p)
	if v == "" {
		return "", nil
	}
	switch o {
	case "U":
		return strings.ToUpper(v), nil
	case "u":
		return strings.ToUpper(v[:1]) + v[1:], nil
	case "L":
		return strings.ToLower(v), nil
	case "l":
		return strings.ToLower(v[:1]) + v[1:], nil
	}
	return "", errSyntax("${parameter@operator}", errors.New("zos: invalid env operator"))
}
