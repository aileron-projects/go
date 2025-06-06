---
title: "TCPサーバ"
linkTitle: "TCPサーバ"
type: docs
weight: 50
categories: []
tags: []
description: ""
---

## 概要

TCPサーバは、TCP(レイヤー4)のサーバを実装することを目的としています。
TCPサーバは、[Goの標準ライブラリ](https://pkg.go.dev/std)では提供されていません。
本機能では[net/httpパッケージ](https://pkg.go.dev/net/http)の[Server](https://pkg.go.dev/net/http#Server)に似たインターフェースでTCPサーバを利用できることを目指しています。

## 機能

TCPサーバが保有する機能は以下の通りです。

### 1. TCPサーバ機能

TCPサーバとして、クライアントからのTCP接続を受け付け、ハンドラーに処理を委譲します。
TCPサーバはTLSに対応しています。

- non-TLS TCPサーバ
- TLS TCPサーバ

サーバの終了処理は以下の２種類があります。

- クローズ
  - `Server.Close`を利用した終了処理
  - 最初にTCP接続の受付を停止（TCPリスナーのクローズ）
  - 続いて既存のTCPコネクションを全てクローズ
- シャットダウン
  - `Server.Shutdown`を利用した終了処理
  - 最初にTCP接続の受付を停止（TCPリスナーのクローズ）
  - 続いて既存のTCPコネクションが自然にクローズされるまで待機
  - シャットダウンタイムアウトが発生した際は残存するコネクションを残してシャットダウン処理を終了

TCPサーバのハンドラのインターフェースは以下のように定義され、
各クライアントのコネクションに対してそれぞれ独立したGoroutineで実行されます。

```go
type Handler interface {
	ServeTCP(ctx context.Context, conn net.Conn)
}
```

### 2. TCPサーバランナー

TCPサーバランナーは、グレースフルシャットダウンを簡単に実装できる機能です。
サーバランナーを利用することでシャットダウン処理の実装の手間を省き、安全にサーバをシャットダウンできます。

利用例は[TCPサーバランナー](#tcpサーバランナー)を参照してください。

## セキュリティに関する特記事項

TCPサーバはTLSを利用可能です。
TLS機能はGo言語の標準機能をそのまま利用しています。

その他、TCPサーバとして特別考慮しているセキュリティはありません。

また、[znet](https://pkg.go.dev/github.com/aileron-projects/go/znet)パッケージの機能を用いることで以下のセキュリティ対策が可能です。

- 同時コネクション数制限 (実装例: [コネクション数制限](#コネクション数制限))
- IPホワイトリスト (実装例: [IPホワイトリスト](#ipホワイトリスト))
- IPブラックリスト (実装例: [IPブラックリスト](#ipブラックリスト))

## 性能に関する特記事項

性能面での特記事項はありません。

## 実装例・使い方

### TCPサーバ

最も基本的なTCPサーバの実装は以下の通りです。
TLSは利用していません。
この実装例のハンドラは、受信したTCPデータを標準出力に出力するのみです。

```go
{{% code source="ex_basic/main.go" %}}
```

### TCPサーバランナー

TCPサーバランナーを利用すると、グレースフルシャットダウンを簡単に実現できます。
シャットダウンタイムアウトが発生した際には、既存のTCPコネクションのクローズ処理も実行してくれます。

TCPサーバランナーの実装例は以下の通りです。
`syscall`パッケージの利用が制限されているプラットフォームもある点に注意してください。

```go
{{% snippet source="ex_runner/main.go" id="main" %}}
```

### TLS

TLSを利用したTCPサーバの実装例は以下の通りです。
TLSサーバのcertファイルとkeyファイルのパスを指定するのみでも起動可能、より細かいTLSの設定をする場合は`Server.TLSConfig`を利用します。

```go
{{% snippet source="ex_tls/main.go" id="main" %}}
```

### コネクション数制限

[znet](https://pkg.go.dev/github.com/aileron-projects/go/znet)パッケージの機能を利用して、
同時コネクション数を制限することが可能です。
実装例は以下となります。一部エラー処理は省略しています。

```go
{{% snippet source="ex_concurrency/main.go" id="main" %}}
```

### IPホワイトリスト

[znet](https://pkg.go.dev/github.com/aileron-projects/go/znet)パッケージの機能を利用して、
ホワイトリスト形式で接続元IPを制限できます。
ホワイトリストは[トライ木](https://en.wikipedia.org/wiki/Trie)で実装されているため、IPアドレスの数が増えても処理時間の劣化はほとんどありません。
許可されないIPアドレスからの接続があった場合、TCPサーバは直ちに当該のTCPコネクションをクローズします。

なお、[znet](https://pkg.go.dev/github.com/aileron-projects/go/znet)パッケージを用いると、ホワイトリストの中でも一部だけは拒否する処理も可能です。

実装例は以下となります。一部エラー処理は省略しています。
この例では、ローカルホストである`127.0.0.1`と`::1`のIPアドレスのみ許可しています。

```go
{{% snippet source="ex_whitelist/main.go" id="main" %}}
```

### IPブラックリスト

[znet](https://pkg.go.dev/github.com/aileron-projects/go/znet)パッケージの機能を利用して、
ブラックリスト形式で接続元IPを制限できます。
ブラックリストは[トライ木](https://en.wikipedia.org/wiki/Trie)で実装されているため、IPアドレスの数が増えても処理時間の劣化はほとんどありません。
許可されないIPアドレスからの接続があった場合、TCPサーバは直ちに当該のTCPコネクションをクローズします。

なお、[znet](https://pkg.go.dev/github.com/aileron-projects/go/znet)パッケージを用いると、ブラックリストの中でも一部だけは許可する処理も可能です。

実装例は以下となります。一部エラー処理は省略しています。
この例では、ローカルホストである`192.168.0.0/16`の範囲にあるIPアドレスを拒否しています。

```go
{{% snippet source="ex_blacklist/main.go" id="main" %}}
```

### Unixドメインソケット

TCPサーバは[Unixドメインソケット](https://en.wikipedia.org/wiki/Unix_domain_socket)を利用できます。

パス名ソケットを利用する場合は、以下のように指定します。

```go
{{% snippet source="ex_socket_path/main.go" id="main" %}}
```

抽象ソケットを利用する場合は、以下のように指定します。

```go
{{% snippet source="ex_socket_abstract/main.go" id="main" %}}
```

## 参考資料
