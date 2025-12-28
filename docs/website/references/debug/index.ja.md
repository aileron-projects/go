# デバッグ

## 概要 {#abstract}

デバッグはデバッグ関連機能を提供します。

デバッグにに関する機能は[zruntime/zdebug](https://pkg.go.dev/github.com/aileron-projects/go/zruntime/zdebug)パッケージにより提供されています。

## 機能 {#features}

### 1. ダンプ出力機能

インスタンスの詳細情報を出力することはデバッグを容易にします。
ダンプ出力機能では、インスタンスの詳細情報を出力できます。

ダンプ出力機能では以下2つの関数が提供されています。
`Dump`関数はビルドタグ`zdebugdump`が指定された時のみ有効となります。
これは、ダンプ出力時には専用のバイナリが必要になる一方でランタイム性能が損なわれません。
つまりデバッグ後にコードを削除する必要がありません。

一方、`DumpAlways`は利用時にビルドタグは不要です。
これは、タグ指定の手間を省き、既存のビルドパイプラインを活用できる一方で、
本番で利用するコードからは削除されなければなりません。

```go
// Build tag `-tags zdebugdump` is required.
zdebug.Dump("dump message", obj1, obj2)

// Tag is not required.
zdebug.DumpAlways("dump message", obj1, obj2)
```

ダンプ機能は内部的に[davecgh/go-spew](https://github.com/davecgh/go-spew)を利用しています。

## セキュリティに関する特記事項 {#security}

予期せず機密情報を出力しないよう、本番環境においてはデバッグ時を除いてダンプ機能を利用しないでください。

## 性能に関する特記事項 {#performance}

性能劣化につながるため、本番環境においてはデバッグ時を除いてダンプ機能を利用しないでください。

## 実装例・使い方 {#example}

### ダンプの基本的な使い方

以下の実装例はプロフィール構造体のダンプを出力します。

```go
--8<-- "references/debug/ex_dump/main.go"
```

上記コードを`-tags zdebugdump`付きで実行した場合、以下のような出力が得られます。

```txt
2025-06-08 00:05:58 [DUMP]dump profile
  | Caller: Pkg=main File=main.go Func=main Line=28
  | ┌── args[0]
  | (*main.profile)(0xc00012e0c0)({
  |  name: (string) (len=8) "john doe",
  |  age: (int) 20,
  |  favorites: ([]string) (len=2 cap=2) {
  |   (string) (len=5) "apple",
  |   (string) (len=6) "orange"
  |  },
  |  experience: (map[string]int) (len=3) {
  |   (string) (len=2) "Go": (int) 3,
  |   (string) (len=3) "C++": (int) 5,
  |   (string) (len=4) "Java": (int) 1
  |  }
  | })
2025-06-08 00:05:58 [DUMP]dump always profile
  | Caller: Pkg=main File=main.go Func=main Line=31
  | ┌── args[0]
  | (*main.profile)(0xc00012e100)({
  |  name: (string) (len=8) "john doe",
  |  age: (int) 20,
  |  favorites: ([]string) (len=2 cap=2) {
  |   (string) (len=5) "apple",
  |   (string) (len=10) "strawberry"
  |  },
  |  experience: (map[string]int) (len=4) {
  |   (string) (len=4) "Java": (int) 1,
  |   (string) (len=4) "Rust": (int) 2,
  |   (string) (len=2) "Go": (int) 3,
  |   (string) (len=1) "C": (int) 6
  |  }
  | })
```

## 参考資料 {#reference}
