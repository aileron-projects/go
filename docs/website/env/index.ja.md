# 環境変数

## 概要

環境変数のパース機能が提供されています。
これらの機能は、YMAL や JSON 等の設定ファイルにおける環境変数の埋め込みなどに利用できます。

環境変数に関する機能は[zos](https://pkg.go.dev/github.com/aileron-projects/go/zos)パッケージにより提供されています。

環境変数名をもとに値を取得する単純なユースケースの場合は、Go 言語標準の[os](https://pkg.go.dev/os)パッケージを利用してください。

## 機能

### 1. 環境変数解決機能

環境変数解決機能は環境変数を値に解決します。
以下のフォーマットの環境変数を解決できます。
これらのルールは基本的に[Shell Parameter Expansion](https://www.gnu.org/software/bash/manual/html_node/Shell-Parameter-Expansion.html)をベースにしています。
parameterは`[0-9a-zA-Z_]+`の文字列です。
wordは`$`を含まない文字列`[^\$]*`です。

ほとんどの場合、本機能は[#2-環境変数置換](#2-環境変数置換)の機能を通して利用されることを想定しています。

1. `${parameter}` : 以下の置換ルール表を参照してください。
2. `${parameter:-word}` : 以下の置換ルール表を参照してください。
3. `${parameter-word}` : 以下の置換ルール表を参照してください。
4. `${parameter:=word}` : 以下の置換ルール表を参照してください。
5. `${parameter=word}` : 以下の置換ルール表を参照してください。
6. `${parameter:?word}` : 以下の置換ルール表を参照してください。
7. `${parameter?word}` : 以下の置換ルール表を参照してください。
8. `${parameter:+word}` : 以下の置換ルール表を参照してください。
9. `${parameter+word}` : 以下の置換ルール表を参照してください。
10. `${parameter:offset}` : `offset`より前の文字を削除します。
11. `${parameter:offset:length}` : `offset`より前と、`offset + length`より後の文字を削除します。
12. `${!prefix*}` : `prefix`を持つパラメータ名を空白区切りで連結します（`${!prefix*}`と同じ動作）。
13. `${!prefix@}` : 現在は #12 と同じ動作です。
14. `${#parameter}` : 値の長さを返します。
15. `${parameter#word}` : 現在は #16 と同じ動作です。
16. `${parameter##word}` : `word` にマッチするプレフィックスを削除します。パターンが指定された場合は最長一致で削除します。
17. `${parameter%word}` : 現在は #18 と同じ動作です。
18. `${parameter%%word}` : `word` にマッチするサフィックスを削除します。パターンが指定された場合は最長一致で削除します。
19. `${parameter/pattern/string}` : 最初に `pattern` にマッチした部分を `string` に置換します。
20. `${parameter//pattern/string}` : `pattern` にマッチしたすべての部分を `string` に置換します。
21. `${parameter/#pattern/string}` : `pattern` に一致した場合、プレフィックスを `string` に置換します。
22. `${parameter/%pattern/string}` : `pattern` に一致した場合、サフィックスを `string` に置換します。
23. `${parameter^pattern}` : `pattern` に一致した場合、最初の文字を大文字に変換します。
24. `${parameter^^pattern}` : `pattern` に一致したすべての文字を大文字に変換します。
25. `${parameter,pattern}` : `pattern` に一致した場合、最初の文字を小文字に変換します。
26. `${parameter,,pattern}` : `pattern` に一致したすべての文字を小文字に変換します。
27. `${parameter@operator}` : `operator` を使って値を処理します。

| #   | expression           | parameter Set and Not Null | parameter Set but Null | parameter Unset |
| --- | -------------------- | -------------------------- | ---------------------- | --------------- |
| 01  | `${parameter}`       | substitute parameter       | substitute null        | substitute null |
| 02  | `${parameter:-word}` | substitute parameter       | substitute word        | substitute word |
| 03  | `${parameter-word}`  | substitute parameter       | substitute null        | substitute word |
| 04  | `${parameter:=word}` | substitute parameter       | substitute word        | assign word     |
| 05  | `${parameter=word}`  | substitute parameter       | substitute null        | assign word     |
| 06  | `${parameter:?word}` | substitute parameter       | error                  | error           |
| 07  | `${parameter?word}`  | substitute parameter       | substitute null        | error           |
| 08  | `${parameter:+word}` | substitute word            | substitute null        | substitute null |
| 09  | `${parameter+word}`  | substitute word            | substitute word        | substitute null |

Cron 式では以下のエイリアスを利用できます。
`CRON_TZ`を除き、各エイリアスはその他の Cron 式と組み合わせて利用することはできません。

| エイリアス名 | 値          | 利用例                  |
| ------------ | ----------- | ----------------------- |
| CRON_TZ      | TZ          | `CRON_TZ=UTC 0 0 * * *` |
| @monthly     | `0 0 1 * *` | `TZ=UTC @monthly`       |
| @weekly      | `0 0 * * 0` | `TZ=UTC @weekly`        |
| @daily       | `0 0 * * *` | `TZ=UTC @daily`         |
| @hourly      | `0 * * * *` | `TZ=UTC @hourly`        |
| @sunday      | `0 0 * * 0` | `TZ=UTC @sunday`        |
| @monday      | `0 0 * * 1` | `TZ=UTC @monday`        |
| @tuesday     | `0 0 * * 2` | `TZ=UTC @tuesday`       |
| @wednesday   | `0 0 * * 3` | `TZ=UTC @wednesday`     |
| @thursday    | `0 0 * * 4` | `TZ=UTC @thursday`      |
| @friday      | `0 0 * * 5` | `TZ=UTC @friday`        |
| @saturday    | `0 0 * * 6` | `TZ=UTC @saturday`      |

`pattern`に指定可能な表記（文字）は以下の通りです。

```text
pattern:
  c       : matches to the character ('$' is not allowed).
  [a-z]   : matches specified character range.
  .*      : matches any length of characters.
  .?      : matches zero or single characters.
```

`operator`として利用可能な文字は以下の通りです。

```text
operator:
  U       : convert all characters to upper case using [strings.ToUpper]
  u       : convert the first character to upper case using [strings.ToUpper]
  L       : convert all characters to lower case using [strings.ToLower]
  l       : convert the first character to lower case using [strings.ToLower]
```

### 2. 環境変数置換

テキストデータ（Goにおける`[]byte`データ）内に含まれる環境変数を解決し、値に置換します。
環境変数の置換に利用できる２つの関数が用意されています。

`EnvSubst`はネスト構造のない環境変数のみ解決できます。
`EnvSubst2`はネスト構造（１階層まで）のある環境変数を置換できます。
つまり、`EnvSubst2`を利用すると`${FOO:-${BAR}}`のような環境変数を置換できます。

```go
func EnvSubst(b []byte) (subst []byte, err error)

func EnvSubst2(b []byte) (subst []byte, err error)
```

### 3. 環境変数ファイルパース機能

環境変数ファイルパース機能は環境変数を定義したファイルから環境変数をパースします。
この機能は`LoadEnv`関数を通して利用します。
環境変数の定義フォーマットは基本的にbashのそれに似せて作られています。

```go
func LoadEnv(b []byte) (map[string]string, error)
```

単一行の変数定義の例を示します。
`>>`の右側にパースされた値を示しています。

```text
FOO=BAR          >> BAR
FOO="BAR"        >> BAR
FOO='BAR'        >> BAR
FOO='B"R'        >> B"R
FOO="B'R"        >> B'R
export FOO=BAR   >> BAR
```

以下のようにクオーテーション（シングルクォーテーションまたはダブルクオーテーション）を利用することで、複数行の値をパースできます。
以下の例の場合、改行（`\n`）はパース時に削除されます。

```text
FOO="
BAR
BAZ
"
```

`#`はコメント行を表します。
行末にコメントを記載する際は`#`の前に半角スペースを必ずおいてください。

```text
# comment
FOO=BAR # comment
```

`\`により文字をエスケープできます。
以下のエスケープパターンをご参照ください。
`>>`の右側にパースされた値を示しています。
`\n`をクォーテーション内部で利用すると改行文字`LF`として扱われます。

```text
FOO=B\"R      >> B"R
FOO=B\'R      >> B'A
FOO="B\"R"    >> B"R
FOO=B\R       >> BR
FOO="B\nR"    >> B<LF>R
```

## セキュリティに関する特記事項

セキュリティに関する特記事項はありません。

## 性能に関する特記事項

性能に関する特記事項はありません。

## 実装例・使い方

### 環境変数の値解決

以下のコードでいくつかの環境変数を解決してみます。

```go
--8<-- "env/ex_envsubst/main.go"
```

上記コードはにより以下の結果が出力されます。

```text
${FOO} ------------ foo
${BAZ:-default} --- default
${BAZ-default}  --- default
${BAZ:=default} --- default
${BAZ=default}  --- default
${BAZ:?default} --- default
${BAZ?default}  --- default
${BAZ:+default} --- default
${BAZ+default}  --- default
${ABC:3} ---------- defg
${ABC:3:3} -------- def
${!ARR*} ---------- ARR_X ARR_Y
${!ARR@} ---------- ARR_X ARR_Y
${#FOO} ---------- 3
${FOO#[a-z]} ----- oo
${FOO##[a-z]} ---- oo
${FOO%[a-z]} ----- fo
${FOO%%[a-z]} ---- fo
${FOO/[a-z]/x} --- xoo
${FOO//[a-z]/x} -- xxx
${FOO/#[a-z]/x} -- xoo
${FOO/%[a-z]/x} -- fox
${FOO^[f]} ------- Foo
${FOO^^[o]} ------ fOO
${BAR,[B]} ------- bAR
${BAR,,[A]} ------ BaR
${FOO@U} --------- FOO
```

### 環境変数の置換

`LoadEnv`による環境変数のパース例を示します。
パースされた値は`os.Setenv`により環境変数としてセットされるとともにマップデータとして返却されます。

```go
--8<-- "env/ex_loadenv/main.go"
```

上記を実行すると以下の値が得られます。
ターミナルへの出力のため、改行など一部エスケープされている点に注意してください。
また、見やすさのために出力をフォーマットしています。

```json
{
  "BAR": "bar",
  "BAZ": "baz",
  "FOO": "foo",
  "MULTILINE_A": "onetwo",
  "MULTILINE_B": "one\ntwo",
  "PASSWORD": "bar",
  "QUOTE_DOUBLE": "double quoted. ' can be used.",
  "QUOTE_DOUBLE_ESCAPE": "double quotation \"escaped\".",
  "QUOTE_SINGLE": "single quoted. \" can be used.",
  "QUOTE_SINGLE_ESCAPE": "single quotation 'escaped'.",
  "SECRET_URL": "http://foo:bar@example.com",
  "URL": "http://example.com",
  "USERNAME": "foo"
}
```

## 参考資料
