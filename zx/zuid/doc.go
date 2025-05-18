// package zid provides some unique ID sources.
//
// Throughout this package,
//
//   - Time precision depends on the [time.Now]
//   - Random value entropy depends on the [rand.Reader]
//
// Random value sources are (cited from [rand.Reader]):
//
//   - On Linux, FreeBSD, Dragonfly, and Solaris, Reader uses getrandom(2).
//   - On legacy Linux (< 3.17), Reader opens /dev/urandom on first use.
//   - On macOS, iOS, and OpenBSD Reader, uses arc4random_buf(3).
//   - On NetBSD, Reader uses the kern.arandom sysctl.
//   - On Windows, Reader uses the ProcessPrng API.
//   - On js/wasm, Reader uses the Web Crypto API.
//   - On wasip1/wasm, Reader uses random_get.
//
// References:
//
//   - RFC 9562 UUIDs: https://datatracker.ietf.org/doc/rfc9562/
//   - RFC 4122 UUID URN Namespace: https://datatracker.ietf.org/doc/rfc4122/
//   - https://github.com/google/uuid
//   - https://github.com/ulid/spec
//   - https://github.com/rs/xid
//   - https://github.com/bwmarrin/snowflake
package zuid
