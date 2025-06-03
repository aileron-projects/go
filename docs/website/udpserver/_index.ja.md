---
title: "UDPサーバ"
linkTitle: "UDPサーバ"
type: docs
weight: 100
categories: []
tags: []
description: ""
---

## 概要

UDPサーバは、UDP(レイヤー4)のサーバを実装することを目的としています。
UDPサーバは、[Goの標準ライブラリ](https://pkg.go.dev/std)では提供されていません。
本機能では[net/httpパッケージ](https://pkg.go.dev/net/http)の[Server](https://pkg.go.dev/net/http#Server)に似たインターフェースでUDPサーバを利用できることを目指しています。

UDPはTCPと異なりコネクションレスのプロトコルです。
そのため、ここで実装されるUDPサーバは仮想コネクションを
仮想コネクション

**用語:**

- UDPサーバ：[UDP(レイヤー４プロトコル)](https://en.wikipedia.org/wiki/User_Datagram_Protocol)のサーバ
- ハンドラー：UDPコネクションを処理する部品
- サーバランナー：UDPサーバを実行するためのヘルパー部品。シャットダウン処理の実装手間を省ける

## 機能

UDPサーバが保有する機能は以下の通りです。

### 1. UDPサーバ機能

UDPサーバとして、クライアントからパケット（UDPデータグラム）を受け取り、ハンドラーに連携します。
UDPサーバはTLSに対応していません。

性能面の理由から、UDPサーバは内部で仮想コネクションを作成します。
仮想コネクションはクライアントの`IP:Port`に対して１つ作成されます。

サーバの終了処理は以下の２種類があります。

- クローズ
  - `Server.Close`を利用した終了処理
  - 最初にUDP接続の受付を停止（UDPリスナーのクローズ）
  - 続いて既存のUDPの処理ループをキャンセル
- シャットダウン
  - `Server.Shutdown`を利用した終了処理
  - 最初にUDP接続の受付を停止（UDPリスナーのクローズ）
  - 続いて既存のUDPの処理ループが自然にクローズされるまで待機
  - シャットダウンタイムアウトが発生した際は残存する処理をそのまま残して終了

UDPサーバのハンドラのインターフェースは以下のように定義され、
各仮想コネクションに対してそれぞれ独立したGoroutineで実行されます。

```go
type Handler interface {
	ServeUDP(ctx context.Context, conn Conn)
}
```

### 2. UDPサーバランナー

UDPサーバランナーは、グレースフルシャットダウンを簡単に実装できる機能です。
サーバランナーを利用することでシャットダウン処理の実装の手間を省き、安全にサーバをシャットダウンできます。

利用例は[UDPサーバランナー](#udpサーバランナー)を参照してください。

## セキュリティに関する特記事項

[znet](https://pkg.go.dev/github.com/aileron-projects/go/znet)パッケージの機能を用いることで以下のセキュリティ対策が可能です。
TCPと異なり、UDPはコネクションレスのプロトコルであるため、TCPサーバにあるような同時接続数の上限設定はできません。

- IPホワイトリスト (実装例: [IPホワイトリスト](#ipホワイトリスト))
- IPブラックリスト (実装例: [IPブラックリスト](#ipブラックリスト))

## 性能に関する特記事項

UDPサーバは、クライアントの`IP:Port`に対して仮想コネクションを作成します。
短時間で多数の`IP:Port`からリクエストを受け取る状況では性能が悪化する可能性があります。

## 実装例・使い方

### UDPサーバ

最も基本的なUDPサーバの実装は以下の通りです。
この実装例のハンドラは、受信したUDPのパケットを標準出力に出力するのみです。

```go
{{% code source="ex_basic/main.go" %}}
```

### UDPサーバランナー

UDPサーバランナーを利用すると、グレースフルシャットダウンを簡単に実現できます。

UDPサーバランナーの実装例は以下の通りです。
`syscall`パッケージの利用が制限されているプラットフォームもある点に注意してください。

```go
{{% snippet source="ex_runner/main.go" id="main" %}}
```

### IPホワイトリスト

[znet](https://pkg.go.dev/github.com/aileron-projects/go/znet)パッケージの機能を利用して、
ホワイトリスト形式で接続元IPを制限できます。
ホワイトリストは[トライ木](https://en.wikipedia.org/wiki/Trie)で実装されているため、IPアドレスの数が増えても処理時間の劣化はほとんどありません。
許可されないIPアドレスからの接続があった場合、UDPサーバは当該パケットを破棄します。

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
許可されないIPアドレスからの接続があった場合、UDPサーバは当該パケットを破棄します。

なお、[znet](https://pkg.go.dev/github.com/aileron-projects/go/znet)パッケージを用いると、ブラックリストの中でも一部だけは許可する処理も可能です。

実装例は以下となります。一部エラー処理は省略しています。
この例では、ローカルホストである`192.168.0.0/16`の範囲にあるIPアドレスを拒否しています。

```go
{{% snippet source="ex_blacklist/main.go" id="main" %}}
```

### Unixドメインソケット

UDPサーバは[Unixドメインソケット](https://en.wikipedia.org/wiki/Unix_domain_socket)を利用できます。

パス名ソケットを利用する場合は、以下のように指定します。

```go
{{% snippet source="ex_socket_path/main.go" id="main" %}}
```

抽象ソケットを利用する場合は、以下のように指定します。

```go
{{% snippet source="ex_socket_abstract/main.go" id="main" %}}
```
