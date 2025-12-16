# ロギング

## 概要 {#abstract}

ログ関連機能を提供します。

ログに関する機能は[zlog](https://pkg.go.dev/github.com/aileron-projects/go/zlog)パッケージにより提供されています。

## 機能 {#features}

### 1. ログ出力機能

ログ出力機能は、ログを出力する機能です。
[zlog](https://pkg.go.dev/github.com/aileron-projects/go/zlog)パッケージでは以下のようにロガーのインターフェースが定められています。
これは特定のログパッケージに依存しないインターフェースとなっています。

```go
type Logger interface {
	DebugEnabled(ctx context.Context) bool
	InfoEnabled(ctx context.Context) bool
	WarnEnabled(ctx context.Context) bool
	ErrorEnabled(ctx context.Context) bool
	DebugContext(ctx context.Context, msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)
}
```

[zlog/zslog](https://pkg.go.dev/github.com/aileron-projects/go/zlog)パッケージは、Go 標準パッケージの[log/slog.Logger](https://pkg.go.dev/log/slog#Logger)にたいしてロガーインターフェースを実装する機能を提供しています。

以下はデフォルトの`slog.Logger`に対して上記の Logger インターフェースを実装する例です。

```go
--8<-- "logging/ex_logger/main.go"
```

### 2. コンテキスト属性操作機能

コンテキスト属性操作機能は`context.Context`にログの属性を保存・取得する機能です。
これによりロガーはコンテキストに紐づく属性をログに含めることが可能です。

属性をコンテキストに保存・取得するために[ContextWithAttrs](https://pkg.go.dev/github.com/aileron-projects/go/zlog#ContextWithAttrs)と[AttrsFromContext](https://pkg.go.dev/github.com/aileron-projects/go/zlog#AttrsFromContext)の２つの関数が提供されています。

```go
ctx := context.Background()

// Save attributes in the context.
ctx = zlog.ContextWithAttrs(ctx, "key", "value")

// Extract attributes from the context.
attrs := zlog.AttrsFromContext(ctx)
```

なお、[zlog/zslog.New](https://pkg.go.dev/github.com/aileron-projects/go/zlog/zslog#New)で作成したロガーを利用した場合、コンテキストに含まれる属性は自動的にログに出力されます。

以下のコードで動作確認ができます。

```go
--8<-- "logging/ex_context_attr/main.go"
```

### 3. コンテキストログレベル設定機能

コンテキストログレベル設定機能は`context.Context`を利用して、対象となるコンテキストに対してログ出力レベルを設定する機能です。
利用しているロガーのログレベルが Info である場合でも、コンテキストのログレベルを Debug に設定すると、当該コンテキストに紐づけられたログレコードは Debug レベルで出力されます。

この機能はデバッグや、コンテキストに応じたログレベルによるログ出力に利用できます。

以下のコードでは、ロガーのログレベルは`Info`になっています。
しかし、コンテキストは`Debug`レベルを指定しているため、Debug ログが出力されます。

ログレベルをコンテキストに設定するために[ContextWithLevel](https://pkg.go.dev/github.com/aileron-projects/go/zlog/zslog#ContextWithLevel)と[LevelFromContext](https://pkg.go.dev/github.com/aileron-projects/go/zlog/zslog#LevelFromContext)の２つの関数が提供されています。

```go
ctx := context.Background()

// Save log level in the context.
ctx = zslog.ContextWithLevel(ctx, slog.LevelDebug)

// Extract log level from the context.
lv := zslog.LevelFromContext(ctx)
```

なお、[zlog/zslog.New](https://pkg.go.dev/github.com/aileron-projects/go/zlog/zslog#New)で作成したロガーを利用した場合、コンテキストに設定されたログレベルは自動的に認識されるようになります。

以下のコードで動作確認ができます。

```go
--8<-- "logging/ex_context_level/main.go"
```

### 4. 論理ファイル機能

論理ファイル機能は、物理ファイルを仮想的にサイズ無限のファイルとして扱う機能です。

ログのファイル出力においてはログファイルのローテーションや世代管理、ログファイル圧縮などを考慮する必要があります。
論理ファイルにより、ロガー自身がこれら物理ファイルの存在や管理を意識する必要がなくなります。
これは Linux における[論理ボリューム管理(LVM)](<https://en.wikipedia.org/wiki/Logical_Volume_Manager_(Linux)>)に似ています。

論理ファイル機能は、アクティブな1つの物理ファイルと、複数の履歴ファイルを管理します。
履歴ファイルのファイル名は固定値（例えば `app.log`）であり、履歴ファイルにはカウンター値や日時を含みます（例えば `app-20060102-150405.log`）。

履歴ファイルのファイル名に指定できるフォーマット指定子は以下の表のとおりです。
フォーマット指定子が指定されていない場合、ファイル管理のために自動的に`%i`が付与されます。
なお、`MaxAge`によるファイル管理を行う場合、日時を含むフォーマット指定子をファイル名に含む必要があります。

| Format | Value                          | Range         |
| ------ | ------------------------------ | ------------- |
| `%Y`   | YYYY 4 digits year             | 0 <= YYYY     |
| `%M`   | MM 2 digits month              | 1 <= MM <= 12 |
| `%D`   | DD 2 digits day of month       | 1 <= DD <= 31 |
| `%h`   | hh 2 digits hour               | 0 <= hh <= 23 |
| `%m`   | mm 2 digits minute             | 0 <= mm <= 59 |
| `%s`   | ss 2 digits second             | 0 <= ss <= 59 |
| `%u`   | unix second with free digits   | 0 <= unix     |
| `%i`   | index with free digits         | 0 <= index    |
| `%H`   | hostname                       |               |
| `%U`   | user id. "-1" on windows       |               |
| `%G`   | user group id. "-1" on windows |               |
| `%p`   | pid (process id)               |               |
| `%P`   | ppid (parent process id)       |               |

論理ファイル機能は以下の項目で物理ファイルを管理します。
これらは組み合わせて利用することも可能です。

- `MaxAge`： 指定時間より古いファイルを削除 (`%Y`, `%M`, `%D` の全て、または `%u` が必須)
- `MaxHistory`： 履歴ファイルが指定数より多い分を削除
- `MaxTotal`： ファイルの合計サイズが指定値を超えた場合に削除

また、履歴ファイルはGzip圧縮により圧縮することも可能です。

## セキュリティに関する特記事項 {#security}

ログ関連機能ではログマスク機能を持ちません。
必要に応じてユーザ側で実装してください。

## 性能に関する特記事項 {#performance}

ログ出力は必ずしも`if`で囲む必要はありません。
可読性と性能を考慮して使い分けることを推奨します。
`Debug`レベルのログはほとんどの場合、本番環境では出力されないのにくわえて出力項目も多くなりがちのため if 文で囲むことが推奨されます。
一方で、`Error`レベルのログはほとんどの場合、本番環境で出力されるため if 文の利用は冗長です。

```go
lg.InfoContext(ctx, "log message")
```

```go
if lg.InfoEnabled(ctx) {
	lg.InfoContext(ctx, "log message")
}
```

## 実装例・使い方 {#example}

### MaxAgeによるファイル管理

以下の実装例は`MaxAge`により履歴管理をおこないます。
１つのファイルのサイズが500バイトを超えないように物理ファイルがローテーションされ、
30秒前より古いファイルは削除されます。

履歴ファイルのファイル名に含まれる`%u`はUnix秒を表しています。
なお、`MaxAge`を利用する場合は、最低限`%Y`/`%M`/`%D`のすべてまたは`%u`を履歴ファイル名に含む必要があります。
時間指定子(`%h`/`%m`/`%s`)が含まれない場合、それらは値がゼロ(00時/00分/00秒)として扱われます。

```go hl_lines="13"
--8<-- "logging/ex_logical_maxage/main.go"
```

### MaxHistoryによるファイル管理

以下の実装例は`MaxHistory`により履歴管理をおこないます。
１つのファイルのサイズが500バイトを超えないように物理ファイルがローテーションされます。
履歴ファイルの数が５つより多くならないように、古いファイルを削除します。

```go hl_lines="13"
--8<-- "logging/ex_logical_maxhistory/main.go"
```

### MaxTotalによるファイル管理

以下の実装例は`MaxTotalBytes`により履歴管理をおこないます。
１つのファイルのサイズが500バイトを超えないように物理ファイルがローテーションされます。
全ての履歴ファイルのファイルサイズの合計が数が2000バイトより大きくならないように、古いファイルから削除します。

```go hl_lines="13"
--8<-- "logging/ex_logical_maxtotal/main.go"
```

### ログのファイル出力

論理ファイルは[io.Writer](https://pkg.go.dev/io#Writer)のインターフェースを実装しているため、
ロガーの出力先として利用することが可能です。

Goの標準パッケージである[log/slog.Handler](https://pkg.go.dev/log/slog#Handler)と組み合わせて利用する場合の利用例を以下に示します。
このような実装により、履歴管理機能付きのロギングを実現することが可能です。

```go
--8<-- "logging/ex_logging/main.go"
```

## 参考資料 {#reference}
