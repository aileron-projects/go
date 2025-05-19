# Go standard library extensions

**This repository provides extensional APIs for the [Go standard library](https://pkg.go.dev/std)**.

<div align="center">

[![GoDoc](https://godoc.org/github.com/aileron-projects/go?status.svg)](http://godoc.org/github.com/aileron-projects/go)
[![AI Docs](https://img.shields.io/badge/AI%20Docs-DeepWiki-blue.svg)](https://deepwiki.com/aileron-projects/go)
[![Go Report Card](https://goreportcard.com/badge/github.com/aileron-projects/go)](https://goreportcard.com/report/github.com/aileron-projects/go)
[![License](https://img.shields.io/badge/License-Apache%202.0-yellow.svg)](./LICENSE)

[![Codecov](https://codecov.io/gh/aileron-projects/go/branch/main/graph/badge.svg?token=L62XLZNFLE)](https://codecov.io/gh/aileron-projects/go)
[![Test Suite](https://github.com/aileron-projects/go/actions/workflows/go-test.yaml/badge.svg?branch=main)](https://github.com/aileron-projects/go/actions/workflows/go-test.yaml?query=branch%3Amain)
[![Check Suite](https://github.com/aileron-projects/go/actions/workflows/check-suite.yaml/badge.svg?branch=main)](https://github.com/aileron-projects/go/actions/workflows/check-suite.yaml?query=branch%3Amain)
[![OpenSourceInsight](https://badgen.net/badge/open%2Fsource%2F/insight/cyan)](https://deps.dev/go/github.com%2Faileron-projects%2Fgo)
[![OSS Insight](https://badgen.net/badge/OSS/Insight/orange)](https://ossinsight.io/analyze/aileron-projects/go)

</div>

## Key Features

- Logging [zlog](./zlog/).
- Debugging [zruntime/zdebug](./zruntime/zdebug/).
- Environmental Variables [zio](./zio/).
- HTTP Reverse Proxy [znet/zhttp](./znet/zhttp/).
- TCP Proxy [znet/ztcp](./znet/ztcp/).
- UDP Proxy [znet/zudp](./znet/zudp/).
- Crontab, Cron Job [ztime/zcron](./ztime/zcron/).
- Rate Limiting [ztime/zrate](./ztime/zrate/).
- Load Balancer [zx/zlb](./zx/zlb/).

## Package Dependency Policy

Package structure, or directory structure, basically follows the [Go standard library](https://pkg.go.dev/std).

All packages in this repository are allowed to use

- [standard packages](https://pkg.go.dev/std)
- [golang.org/x](https://pkg.go.dev/golang.org/x)
- Third party packages
  - [github.com/kr/pretty](https://pkg.go.dev/github.com/kr/pretty)
  - [github.com/davecgh/go-spew/spew](https://pkg.go.dev/github.com/davecgh/go-spew/spew)

A package can contain package of higher-level APIs in its subdirectories.
Higher level APIs can use lower level APIs.
That means a package can use parent packages and cannot use child packages.

For example, in the following package structure,

- package `lowapi` cannot use neither `middleapi` nor `highapi`
- package `middleapi` can use `lowapi` and cannot use `highapi`
- package `highapi` can use both `lowapi` and `middleapi`

```text
lowapi/  <────────┐  <──┐
│                 |     |
└── middleapi/  ──┘  <──┤
    │                   |
    └── highapi/  ──────┘
```

## Tested Environment

Operating System:

- `Linux` ([ubuntu-latest](https://github.com/actions/runner-images))
- `Windows` ([windows-latest](https://github.com/actions/runner-images))
- `macOS` ([macos-latest](https://github.com/actions/runner-images))

Go:

- Current Stable (Latest patch version of `go 1.(N).x`)
- Previous Stable (Latest patch version of `go 1.(N-1).x`)
- Minimum Requirement (Currently `go 1.24.0`)

Where `N` is the current latest minor version.
See the Go official release page [Stable versions](https://go.dev/dl/).

In addition to the environment above, following platforms are tested on ubuntu
using [QEMU User space emulator](https://www.qemu.org/docs/master/user/main.html).

- `amd64`
- `arm/v5`
- `arm/v6`
- `arm/v7`
- `arm64`
- `ppc64`
- `ppc64le`
- `riscv64`
- `s390x`
- `loong64`
- `386`
- `mips`
- `mips64`
- `mips64le`
- `mipsle`

## Release Cycle

- Releases are made as needed.
- Versions follow [Semantic Versioning](https://semver.org/).
  - `vZ.Y.Z`
  - `vZ.Y.Z-rc.N`
  - `vZ.Y.Z-beta.N`
  - `vZ.Y.Z-alpha.N`

## License

Apache 2.0
