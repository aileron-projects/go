# レートリミット

## 概要

レートリミットは、任意の処理に対して実行回数の制限を適用します。
レートリミットが利用される典型的な例は、APIのリクエスト数の制御です。
APIのレートリミットとして最も一般的に用いられるのは[TokenBucketアルゴリズム](#3-tokenbucket)によるレートリミットです。

レートリミットは[ztime/zrate](https://pkg.go.dev/github.com/aileron-projects/go/ztime/zrate)パッケージにより提供されています。

## 機能

### 1. LeakyBucket

本機能は[LeakyBucketアルゴリズム](https://en.wikipedia.org/wiki/Leaky_bucket)によるレートリミットを実現します。

LeakyBucketのリミッター作成時に指定できるパラメータは以下の通りです。

1. `queueSize`: キューサイズ
2. `interval`: リクエストを処理する間隔（デキューの間隔）

アルゴリズムの概要は以下になります。

- サイズ`N`のキューを作成する。
  - キューはFIFOで処理される。
- リクエストはキューに挿入される。
  - キューが満杯の場合、リクエストを拒否する。
- 一定のインターバルでキューの先頭からリクエストを取り出し処理する。

LeakyBucketは以下の2つのメソッドを持ちます。
`AllowNow`はキューが空の場合にのみ、有効なトークンを返します。
従って、LeakyBucketを利用する場合は基本的には`WaitNow`を利用し、リクエストが処理されるのを待機します。
`Accept`は`WaitNow`のエイリアスです。

```go
func AllowNow() Token
func WaitNow(ctx context.Context) Token
func Accept(ctx context.Context) Token
```

APIのレートリミットに使用する例は、[LeakyBucketによるAPIレートリミット](#leakybucketによるapiレートリミット)を参照してください。

### 2. MaxConcurrency

本機能はMaxConcurrencyアルゴリズムによるレートリミットを実現します。

Concurrencyのリミッター作成時に指定できるパラメータは以下の通り。

1. `concurrency`: 最大同時実行数

アルゴリズムの概要は以下になります。

- `N`個のセマフォ変数を作成する。
- リクエストは1つのセマフォ変数をロックし、処理を開始する。
  - セマフォ変数が全てロックされている場合、直ちに処理をエラーとするか、ロックを獲得できるまで待機する。
- 処理が完了するとセマフォ変数のロックを開放する。

Concurrencyリミッターは以下の2つのメソッドを持ちます。
`AllowNow`はセマフォ変数が全てロックされている場合は直ちに無効なトークンを返却、一方`WaitNow`はセマフォ変数のロックを獲得できるまで待機する。
`Accept`は`AllowNow`のエイリアスです。

```go
func AllowNow() Token
func WaitNow(ctx context.Context) Token
func Accept(ctx context.Context) Token
```

APIのレートリミットに使用する例は、[ConcurrencyによるAPIレートリミット](#concurrencyによるapiレートリミット)を参照してください。

### 3. TokenBucket

本機能は[TokenBucketアルゴリズム](https://en.wikipedia.org/wiki/Token_bucket)によるレートリミットを実現します。
TokenBucketはAPIのレートリミットにおいて広く用いられているアルゴリズムです。
TokenBucketアルゴリズムはバースト（瞬間的な処理数の増加）を許容するアルゴリズムです。

TokenBucketのリミッター作成時に指定できるパラメータは以下の通りです。

1. `bucketSize`: バケットサイズ
2. `fillRate`: トークン補充レート（基本的には1秒間割合を指す）

アルゴリズムの概要は以下になります。

- サイズ`N`のトークン用バケットを作成する。
  - バケットには一定間隔で`r`個のバケットが補充される。
  - バケットが満杯であれば、それ以上トークンは補充されない。
- リクエストはバケットから1つのトークンを消費して処理を開始する。
  - トークンが存在しなければ処理を開始できない。トークンの補充を待機するか直ちに処理を中断する。

TokenBucketリミッターは以下の2つのメソッドを持ちます。
`AllowNow`は1つのトークンを消費して処理を開始します。トークンが存在しなければ直ちに無効なトークンが返却されるため、処理は中断されます。
一方、`WaitNow`は同様にトークンを消費して処理を開始しますが、トークンが存在しなければ補充されるまで待機します。
`Accept`は`AllowNow`のエイリアスです。

```go
func AllowNow() Token
func WaitNow(ctx context.Context) Token
func Accept(ctx context.Context) Token
```

APIのレートリミットに使用する例は、[TokenBucketによるAPIレートリミット](#tokenbucketによるapiレートリミット)を参照してください。

### 4. FixedWindow

本機能はFixedWindowアルゴリズムによるレートリミットを実現します。
FixedWindowは一定区間における処理数を厳密に制限できる一方、異なる区間の境界において制限値の2倍の処理が実行される可能性のあるアルゴリズムです。

FixedWindowのリミッター作成時に指定できるパラメータは以下の通りです。

1. `limit`: 処理数上限

アルゴリズムの概要は以下になります。

- サイズ`N`のトークン用バケットを作成する。
  - バケットには一定間隔で満杯になるようにトークンが補充される。
- リクエストはバケットから1つのトークンを消費して処理を開始する。
  - トークンが存在しなければ処理を開始できない。トークンの補充を待機するか直ちに処理を中断する。

FixedWindowリミッターは以下の2つのメソッドを持ちます。
`AllowNow`は1つのトークンを消費して処理を開始します。トークンが存在しなければ直ちに無効なトークンが返却されるため、処理は中断されます。
一方、`WaitNow`は同様にトークンを消費して処理を開始しますが、トークンが存在しなければ補充されるまで待機します。
`Accept`は`AllowNow`のエイリアスです。

```go
func AllowNow() Token
func WaitNow(ctx context.Context) Token
func Accept(ctx context.Context) Token
```

APIのレートリミットに使用する例は、[FixedWindowによるAPIレートリミット](#fixedwindowによるapiレートリミット)を参照してください。

## セキュリティに関する特記事項

セキュリティに関する特記事項はありません。

## 性能に関する特記事項

性能に関する特記事項はありません。

## 実装例・使い方

### LeakyBucketによるAPIレートリミット

以下の例は、LeakyBucketアルゴリズムによるAPIレートリミットの簡単な実装例です。
必要に応じてパス単位などでレートリミットを適用することも可能です。

以下の例では、キューサイズが10、リクエストを処理するインターバルは100ミリ秒に設定されています。
また、リクエストが受理された場合、0-100ミリ秒間の間でランダムに待機した後200ステータスを返却します。

```go
--8<-- "rate/ex_leakybucket/main.go"
```

### ConcurrencyによるAPIレートリミット

ConcurrencyアルゴリズムによるAPIレートリミットの簡単な実装例です。
必要に応じてパス単位などでレートリミットを適用することも可能です。

以下の例では、同時処理数が10に設定されています。
また、リクエストが受理された場合、0-100ミリ秒間の間でランダムに待機した後200ステータスを返却します。

```go
--8<-- "rate/ex_concurrent/main.go"
```

### TokenBucketによるAPIレートリミット

TokenBucketアルゴリズムによるAPIレートリミットの簡単な実装例です。
必要に応じてパス単位などでレートリミットを適用することも可能です。

以下の例では、バケットサイズが10、トークン補充割合が10/秒に設定されています。
また、リクエストが受理された場合、0-100ミリ秒間の間でランダムに待機した後200ステータスを返却します。

```go
--8<-- "rate/ex_tokenbucket/main.go"
```

### FixedWindowによるAPIレートリミット

FixedWindowアルゴリズムによるAPIレートリミットの簡単な実装例です。
必要に応じてパス単位などでレートリミットを適用することも可能です。

以下の例では、上限が10リクエスト/秒に設定されています。
また、リクエストが受理された場合、0-100ミリ秒間の間でランダムに待機した後200ステータスを返却します。

```go
--8<-- "rate/ex_fixedwindow/main.go"
```

### SlidingWindowによるAPIレートリミット

SlidingWindowアルゴリズムによるAPIレートリミットの簡単な実装例です。
必要に応じてパス単位などでレートリミットを適用することも可能です。

以下の例では、上限が10に設定されています。
また、リクエストが受理された場合、0-100ミリ秒間の間でランダムに待機した後200ステータスを返却します。

```go
--8<-- "rate/ex_slidingwindow/main.go"
```

## 参考資料
