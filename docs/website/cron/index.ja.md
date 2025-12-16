# Cron/Cronジョブ {#cron-cronjob}

## 概要 {#abstract}

Cronはスケジューリングに関する機能で、ジョブの定期実行などで利用できます。
この機能は、例えばログローテーションなどに活用できます。

Cronに関する機能は[ztime/zcron](https://pkg.go.dev/github.com/aileron-projects/go/ztime/zcron)パッケージにより提供されています。

## 機能 {#features}

### 1. Cron式パース機能 {#cron-parse}

本機能は[Cron式](https://en.wikipedia.org/wiki/Cron)をパースする機能です。
Cron式には一般的な標準仕様が存在しないため、本機能は以下の仕様に従ってパースします。

Cron式のフォーマットは以下の通りです。
タイムゾーン、および秒フィールドの指定はオプショナルです。

```txt
TZ=UTC * * * * * *
|      | | | | | |
|      | | | | | |- Day of week
|      | | | | |--- Month
|      | | | |----- Day of month
|      | | |------- Hour
|      | |--------- Minute
|      |----------- Second (Optional)
|------------------ Timezone (Optional)
```

それぞれのフィールドは表に示した値を利用できます。
`Day of month`と`Day of week`はAND条件で評価されます。

| フィールド名 | 指定必須 | 値              | 特殊文字        |
| ------------ | -------- | --------------- | --------------- |
| Timezone     | No       | Timezone name   |                 |
| Second       | No       | 0-59            | `*` `/` `,` `-` |
| Minute       | Yes      | 0-59            | `*` `/` `,` `-` |
| Hours        | Yes      | 0-23            | `*` `/` `,` `-` |
| Day of month | Yes      | 1-31            | `*` `/` `,` `-` |
| Month        | Yes      | 1-12 or JAN-DEC | `*` `/` `,` `-` |
| Day of week  | Yes      | 0-6 or SUN-SAT  | `*` `/` `,` `-` |

**エイリアス:**

Cron式では以下のエイリアスを利用できます。
`CRON_TZ`を除き、各エイリアスはその他のCron式と組み合わせて利用することはできません。

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

**every表記:**

利便性のために、`@every`表記が用意されています。
この表記は1秒以上、24時間未満で定期的に実行する場合に利用できます。

なおパース可能な表記はGo標準パッケージの[time.Parse](https://pkg.go.dev/time#Parse)をご参照ください。

| Duration  | Resolved Cron          | Notes                 |
| --------- | ---------------------- | --------------------- |
| -1s       | ERROR                  | Duration must be >0   |
| 0s        | ERROR                  | Duration must be >0   |
| 1s        | `*/1 * * * * *`        |                       |
| 1m        | `0 */1 * * * *`        |                       |
| 1h        | `0 0 */1 * * *`        |                       |
| 61s       | `*/1 */1 * * *`        |                       |
| 15m30s    | `*/30 */15 * * * *`    |                       |
| 65m30s    | `*/30 */5 */1 * * *`   |                       |
| 1h30m     | `0 */30 */1 * * *`     |                       |
| 23h59m59s | `*/59 */59 */23 * * *` |                       |
| 24h       | ERROR                  | Duration must be <24h |

### 2. スケジュール時刻決定機能 {#cron-scheduling}

スケジュール時刻決定機能は、パースされたCron式によりスケジュールされる時刻を評価する機能です。
スケジュール時刻を取得する関数として以下の2つが用意されています。

`Next`は現在時刻以降で最初にスケジュールされる時刻を返します。
`NextAfter`は指定されて時刻以降で最初にスケジュールされる時刻を返します。
なお、`NextAfter`により返却される時刻のタイムゾーンは引数として受け取ったタイムゾーンになります。

基本的な利用方法は[cronの基本的な利用](#basic-usage)を参照ください。

```go
func Next() time.Time
func NextAfter(t time.Time) time.Time
```

### 3. ジョブ実行機能 {#job-execution}

ジョブ実行機能はCronによりスケジュールされた時刻にジョブを実行する機能です。
リスケジューリングなどの高度な機能は実装されていません。
デフォルトでは以下２つの項目を制御可能です。

- 最大同時実行数: ある時刻で同時に起動しているジョブの最大数
- 最大リトライ回数: ジョブが失敗した際にリトライする最大回数

基本的な利用方法は[ジョブ実行例](#job-example)を参照ください。

## セキュリティに関する特記事項 {#security}

セキュリティに関する特記事項はありません。

## 性能に関する特記事項 {#performance}

性能に関する特記事項はありません。

## 実装例・使い方 {#example}

### Cronの基本的な利用 {#basic-usage}

cronの基本的な利用方法は以下のようになります。

```go
--8<-- "cron/ex_basic/main.go"
```

### ジョブ実行例 {#job-example}

ジョブの実行例は以下の通りです。
この例では3秒ごとに時刻を表示しています。
`Start`の実行ではGoroutineを利用していないので、`Run`を実行している行で処理はブロックされます。
ジョブを非同期で実行したい場合は`go cron.Start()`のように実装します。

```go
--8<-- "cron/ex_job/main.go"
```

## 参考資料 {#reference}

- [wiki/Cron](https://en.wikipedia.org/wiki/Cron)
- [crontab(5) - Linux manual page](https://man7.org/linux/man-pages/man5/crontab.5.html)
- [go-co-op/gocron](https://github.com/go-co-op/gocron)
- [jasonlvhit/gocron](https://github.com/jasonlvhit/gocron)
- [crontab guru](https://crontab.guru/)
