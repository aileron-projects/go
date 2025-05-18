package zos

import (
	"bytes"
	"errors"
	"os"
)

// LoadEnv loads environmental variable from the given bytes.
// Typically LoadEnv is used for loading .env file.
// LoadEnv resolves embedded environmental variables in the b.
// THe syntax for substituting environmental variable follows the
// specification of [ResolveEnv].
//
// References:
//   - https://github.com/joho/godotenv
//   - https://github.com/motdotla/dotenv
//
// Input specifications:
//
//	Single line:
//		# Following declaration results in BAR.
//		# Single quotes and double quotes are removed if entire value is enclosed.
//		# "export" can be placed before env name.
//		FOO=BAR          >> BAR
//		FOO="BAR"        >> BAR
//		FOO='BAR'        >> BAR
//		FOO='B"R'        >> B"R
//		FOO="B'R"        >> B'R
//		export FOO=BAR   >> BAR
//
//	Multiple line:
//		# The following definition of FOO results in "BARBAZ".
//		# Line breaks of LF or CRLF are removed.
//		# BOTH single quotes and double quotes can be used to enclose multiple lines.
//		FOO="
//		BAR
//		BAZ
//		"
//
//	Comments:
//		# Sharp '#' can be used for commenting.
//		# It must not be in the scope of single or double quotes.
//		# It must have at least 1 white space before '#' if the comment is after value.
//		# comment            >> Comment is appropriately parsed.
//		FOO=BAR # comment    >> Comment is appropriately parsed.
//		FOO=BAR# comment     >> '#' is not parsed as comment. It considered as a part of value.
//
//	Escapes:
//		# '\\' can be used for escaping character following the 3 rules.
//		# 1. '\\' always escapes special character of ', ", \\, #
//		# 2. '\\' is ignored when it is not in the scope of single or double quotes.
//		# 3. '\\'n or "\n" in the scope of single or doubles quotes results in line breaks of LF.
//		FOO=B\"R      >> B"R
//		FOO=B\'R      >> B'A
//		FOO="B\"R"    >> B"R
//		FOO=B\R       >> BR (Its not in a scope of single or double quotes.)
//		FOO="B\nR"    >> B<LF>R (\n is, if in a scope of quotes, converted into a line break.)
//
//	Environmental variables:
//		# LoadEnv resolves environmental variables.
//		FOO=BAR${BAZ}
func LoadEnv(b []byte) (map[string]string, error) {
	envs := map[string]string{}
	inSingleQuote := false
	inDoubleQuote := false
	multilineName := ""
	multilineValue := ""
	bb := b
	var found bool
	var line []byte
	for {
		line, bb, found = bytes.Cut(bb, []byte("\n"))
		if len(line) == 0 && !found {
			break // End of file.
		}
		line = bytes.Trim(line, "\t\n\f\r ")
		line, err := EnvSubst(line) // Replace environmental variable if exists.
		if err != nil {
			return nil, &EnvError{Err: err, Type: typeLoad}
		}

		if inSingleQuote || inDoubleQuote {
			var val string
			val, inSingleQuote, inDoubleQuote = scanValue(line, inSingleQuote, inDoubleQuote)
			multilineValue += val
			if !inSingleQuote && !inDoubleQuote {
				envs[multilineName] = multilineValue
				_ = os.Setenv(multilineName, multilineValue)
				multilineName = ""  // Reset variable.
				multilineValue = "" // Reset variable.
			}
		} else {
			line = bytes.TrimPrefix(line, []byte("export"))
			line = bytes.TrimLeft(line, "\t ")
			name, rest, err := scanName(line)
			if err != nil {
				return nil, &EnvError{Err: err, Type: typeLoad}
			}
			if name == "" {
				continue // Maybe comment line.
			}
			var val string
			val, inSingleQuote, inDoubleQuote = scanValue(rest, inSingleQuote, inDoubleQuote)
			if inSingleQuote || inDoubleQuote {
				multilineName = name
				multilineValue = val
			} else {
				envs[name] = val
				_ = os.Setenv(name, val)
			}
		}
	}
	if inSingleQuote || inDoubleQuote {
		return nil, &EnvError{Type: typeLoad, Info: "Quotation is not closed in variable " + multilineName}
	}
	return envs, nil
}

// scanName scans the given line and looks for environmental variable name.
// If the line is comment, it returns empty name.
// If a variable name found, it returns the name and the rest of the line after '='.
//
// Example input patterns:
//   - FOO=bar
//   - FOO="bar"
//   - FOO="bar" #comment
//   - FOO="bar
//   - #comment
func scanName(b []byte) (name string, rest []byte, err error) {
	if len(b) == 0 || b[0] == '#' {
		return "", nil, nil
	}
	if b[0] == '=' {
		return "", nil, errors.New("zos: variable name not found")
	}
	for i, c := range b {
		switch {
		case '0' <= c && c <= '9':
		case 'a' <= c && c <= 'z':
		case 'A' <= c && c <= 'Z':
		case c == '_':
		case c == '=':
			return string(b[:i]), b[i+1:], nil
		default:
			return "", nil, errors.New("zos: character not allowed for variable name `" + string(c) + "`")
		}
	}
	return "", nil, errors.New("zos: invalid line")
}

// scanValue scans value line.
// It returns the value and the flag of
// in-single-quote or in-double-quote.
// isq and idq don't become true simultaneously.
//
// Example input patterns:
//
//	foobar        <-- Non quoted value
//	"foobar"      <-- Quoted value
//	"foobar       <-- Quoted value. Double quote is not closed.
//	'foobar       <-- Quoted value. Single quote is not closed.
//	"foo"'bar     <-- Quoted value. Single quote is not closed.
//	"foo"'bar'    <-- Non quoted value. All quotes are closed.
//	foo#comment   <-- Non quoted value. The '#' is treated as a value.
//	foo #comment  <-- Commented value. Comments are ignored.
//	#comment      <-- Commented value. Comments are ignored.
func scanValue(b []byte, sq, dq bool) (val string, isq, idq bool) {
	v := make([]byte, 0, len(b))
	escaped := false
	inSingleQuote := sq
	inDoubleQuote := dq
	for _, c := range b {
		if escaped {
			if c == '\\' || c == '\'' || c == '"' || c == '#' {
				v = append(v, c)
				escaped = false
				continue
			}
			if (inSingleQuote || inDoubleQuote) && c == 'n' {
				v = append(v, '\n')
				escaped = false
				continue
			}
			v = append(v, c)
			escaped = false
			continue
		}
		switch c {
		case '\\':
			escaped = true
		case '\'':
			if inDoubleQuote {
				v = append(v, c)
			} else {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if inSingleQuote {
				v = append(v, c)
			} else {
				inDoubleQuote = !inDoubleQuote
			}
		case '#':
			if !inSingleQuote && !inDoubleQuote && (len(v) > 1 && v[len(v)-1] == ' ') {
				v = bytes.TrimRight(v, " ")
			}
			return string(v), inSingleQuote, inDoubleQuote
		default:
			v = append(v, c)
		}
	}
	if escaped {
		v = append(v, '\\')
	}
	return string(v), inSingleQuote, inDoubleQuote
}
